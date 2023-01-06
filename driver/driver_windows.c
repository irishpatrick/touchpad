#include "driver.h"

#include <assert.h>
#include <stdint.h>
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#define QUEUE_DEPTH 100

static INPUT input_queue[QUEUE_DEPTH];
static size_t input_queue_size = 0;

static INPUT* input_queue_push_get(void)
{
    INPUT* input = &input_queue[input_queue_size];
    ++input_queue_size;
    return input;
}

int driver_create_device(void)
{
    ZeroMemory(input_queue, sizeof(input_queue));
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
    // TODO implement
    return 0;
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
}