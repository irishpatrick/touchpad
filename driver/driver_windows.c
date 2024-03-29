#include "driver.h"

#include <assert.h>
#include <stdint.h>
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#define QUEUE_DEPTH 100

static INPUT input_queue[QUEUE_DEPTH];
static size_t input_queue_size = 0;

static void start_event_loop_thread(void)
{

}

static void stop_event_loop_thread(void)
{

}

static INPUT* input_queue_push_get(void)
{
    INPUT* input = &input_queue[input_queue_size];
    ++input_queue_size;
    return input;
}

int driver_create_device(void)
{
    ZeroMemory(input_queue, sizeof(input_queue));
    start_event_loop_thread();
    return 0;
}

int driver_mouse_rel(int x, int y)
{
    INPUT* input = input_queue_push_get();

    input->type = INPUT_MOUSE;
    input->mi.mouseData = 0;
    input->mi.dwFlags = MOUSEEVENTF_MOVE;
    input->mi.dwExtraInfo = 0;
    input->mi.time = 0;
    
    input->mi.dx = x;
    input->mi.dy = y;

    return driver_report();
}

int driver_mouse_btn(int left, int middle, int right)
{
    INPUT* input = input_queue_push_get();

    input->type = INPUT_MOUSE;
    input->mi.mouseData = 0;
    input->mi.dwFlags = 
        MOUSEEVENTF_LEFTDOWN & (left == 1) | 
        MOUSEEVENTF_MIDDLEDOWN & (middle == 1) | 
        MOUSEEVENTF_RIGHTDOWN & (right == 1) |
        MOUSEEVENTF_LEFTUP & (left == 0) | 
        MOUSEEVENTF_MIDDLEUP & (middle = 0) |
        MOUSEEVENTF_RIGHTUP & (right == 0);
    input->mi.dwExtraInfo = 0;
    input->mi.time = 0;

    return driver_report();
}

int driver_report(void)
{
    int result = 0;
    UINT num_sent = SendInput(input_queue_size, input_queue, input_queue_size * sizeof(INPUT));
    if (num_sent != input_queue_size)
    {
        result = 1;
    }

    input_queue_size = 0;
    return result;
}

void driver_destroy_device(void)
{
    driver_report();
    stop_event_loop_thread();
}