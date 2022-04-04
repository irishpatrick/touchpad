#include "driver.h"

#include <assert.h>
#include <linux/input.h>
#include <libevdev/libevdev.h>
#include <libevdev/libevdev-uinput.h>
#include <stdlib.h>
#include <errno.h>

static struct libevdev* dev = NULL;
static struct libevdev_uinput* uidev = NULL;

int driver_create_device(void)
{
    int err;
    
    dev = libevdev_new();
    libevdev_set_name(dev, "test device");
    libevdev_enable_event_type(dev, EV_REL);
    libevdev_enable_event_code(dev, EV_REL, REL_X, NULL);
    libevdev_enable_event_code(dev, EV_REL, REL_Y, NULL);
    libevdev_enable_event_type(dev, EV_KEY);
    libevdev_enable_event_code(dev, EV_KEY, BTN_LEFT, NULL);
    libevdev_enable_event_code(dev, EV_KEY, BTN_MIDDLE, NULL);
    libevdev_enable_event_code(dev, EV_KEY, BTN_RIGHT, NULL);

    err = libevdev_uinput_create_from_device(dev, LIBEVDEV_UINPUT_OPEN_MANAGED, &uidev);
    if (err != 0)
    {
        return err;
    }

    return 0;
}

int driver_mouse_rel(int x, int y)
{
    int err;

    err = libevdev_uinput_write_event(uidev, EV_REL, REL_X, x);
    if (err != 0)
    {
        return err;
    }
    
    err = libevdev_uinput_write_event(uidev, EV_KEY, REL_Y, y);
    if (err != 0)
    {
        return err;
    }

    return driver_report();
}

int driver_mouse_btn(int left, int middle, int right)
{
    int err;

    err = libevdev_uinput_write_event(uidev, EV_REL, BTN_LEFT, left);
    if (err != 0)
    {
        return err;
    }
    
    err = libevdev_uinput_write_event(uidev, EV_REL, BTN_MIDDLE, middle);
    if (err != 0)
    {
        return err;
    }
    
    err = libevdev_uinput_write_event(uidev, EV_REL, BTN_RIGHT, right);
    if (err != 0)
    {
        return err;
    }

    return driver_report();
}

int driver_report(void)
{
    int err;

    err = libevdev_uinput_write_event(uidev, EV_SYN, SYN_REPORT, 0);
    if (err != 0)
    {
        return err;
    }

    return 0;
}

void driver_destroy_device(void)
{
    libevdev_uinput_destroy(uidev);
}

