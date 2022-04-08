#include <curl/curl.h>
#include <gtk/gtk.h>
#include <pthread.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

static pthread_t conn_thread;
static pthread_mutex_t conn_mutex;
static bool is_connected;
static bool conn_thread_kill = false;
static char qr_url[256];
static GtkWidget* canvas;

struct conn_thread_data
{
    GtkWidget* darea;
    const char* endpt;
};

static size_t my_write(void* ptr, size_t size, size_t nmemb, FILE* fp)
{
    return fwrite(ptr, size, nmemb, fp);
}

static size_t my_read(void* ptr, size_t size, size_t nmemb, FILE* fp)
{
    return fread(ptr, size, nmemb, fp);
}

static void print_hello (GtkWidget* widget, gpointer data)
{
    g_print("Hello World\n");
}

static gboolean draw_callback(GtkWidget* widget, cairo_t* cr, gpointer data)
{
    guint width;
    guint height;
    GdkRGBA color;
    GtkStyleContext* ctx;
    
    ctx = gtk_widget_get_style_context(widget);
    width = gtk_widget_get_allocated_width(widget);
    height = gtk_widget_get_allocated_height(widget);

    gtk_render_background(ctx, cr, 0, 0, width, height);

    if (is_connected)
    {
        cairo_arc (cr, width / 2.0, height / 2.0, MIN (width, height) / 2.0, 0, 2 * G_PI);
        gtk_style_context_get_color (ctx, gtk_style_context_get_state (ctx), &color);
        cairo_set_source_rgba (cr, 1.0, 0.0, 0.0, 1.0);
        cairo_fill (cr);
    }
    else
    {
        cairo_arc (cr, width / 2.0, height / 2.0, MIN (width, height) / 2.0, 0, 2 * G_PI);
        gtk_style_context_get_color (ctx, gtk_style_context_get_state (ctx), &color);
        cairo_set_source_rgba (cr, 0.0, 0.0, 0.0, 1.0);
        cairo_fill (cr);
    }

    return FALSE;
}

static void* conn_worker(void* ptr)
{
    CURL* curl;

    curl = curl_easy_init();
    if (curl)
    {
        struct conn_thread_data* data = (struct conn_thread_data*)ptr;
        gchar* url = data->endpt;
        FILE* fp = tmpfile();
        char buffer[8192];

        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, fp);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, my_write);
        curl_easy_setopt(curl, CURLOPT_READFUNCTION, my_read);
    
        while (true)
        {
            curl_easy_perform(curl);

            fseek(fp, 0L, SEEK_SET);
            fread(buffer, 1, 8192, fp);
            buffer[8191] = 0;
            if (strlen(buffer) < 1)
            {
                printf("no response!\n");
                sleep(1);
            }
            else
            {
                pthread_mutex_lock(&conn_mutex);
                is_connected = true;
                strncpy(qr_url, buffer, 256);
                pthread_mutex_unlock(&conn_mutex);

                gtk_widget_queue_draw(canvas);

                break;
            }

            pthread_mutex_lock(&conn_mutex);
            if (conn_thread_kill)
            {
                pthread_mutex_unlock(&conn_mutex);

                break;
            }
            pthread_mutex_unlock(&conn_mutex);
        }

        fclose(fp);
        curl_easy_cleanup(curl);
    }

    return NULL;
}

static void activate(GtkApplication* app, gpointer user_data)
{
    GtkWidget* window;
    GtkWidget* button;
    GtkWidget* box;

    // init curl
    curl_global_init(CURL_GLOBAL_ALL);

    // build gui
    window = gtk_application_window_new(app);
    gtk_window_set_title(GTK_WINDOW(window), "Touchpad Frontend");
    gtk_window_set_default_size(GTK_WINDOW(window), 512, 480);

    box = gtk_vbox_new(0, 1);
    
    canvas = gtk_drawing_area_new();
    gtk_widget_set_size_request(canvas, 512, 460);
    g_signal_connect(canvas, "draw", G_CALLBACK(draw_callback), NULL);
    gtk_box_pack_start(GTK_BOX(box), canvas, 1, 1, 1);
    
    button = gtk_button_new_with_label("Hello World");
    g_signal_connect(button, "clicked", G_CALLBACK(print_hello), NULL);
    gtk_box_pack_start(GTK_BOX(box), button, 1, 0, 1);
    
    gtk_container_add(GTK_CONTAINER(window), box);

    gtk_widget_show_all(window);
    gtk_window_present(GTK_WINDOW(window));

    gtk_widget_queue_draw(canvas);
}

int main(int argc, char** argv)
{
    GtkApplication* app;
    int status;
    int ret;

    struct conn_thread_data data;
    data.endpt = "http://localhost:8080/url";
    data.darea = canvas;
    ret = pthread_create(&conn_thread, NULL, conn_worker, (void*)(&data));

    app = gtk_application_new("org.gtk.example", G_APPLICATION_FLAGS_NONE);
    g_signal_connect(app, "activate", G_CALLBACK(activate), NULL);
    status = g_application_run(G_APPLICATION(app), argc, argv);
    g_object_unref(app);

    if (!is_connected)
    {
        pthread_mutex_lock(&conn_mutex);
        conn_thread_kill = true;
        pthread_mutex_unlock(&conn_mutex);
    }
    pthread_join(conn_thread, NULL);

    return status;
}

