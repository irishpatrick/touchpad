#include "qrcodegen.h"
#include "qrcode.h"
#include <assert.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

cairo_surface_t* generate_qr_code(const char* data)
{
    cairo_surface_t* surf = NULL;
    QRCode qrcode;
    int8_t ret;
    uint8_t qrData[qrcode_getBufferSize(3)];

    // generate qr
    ret = qrcode_initText(&qrcode, qrData, 3, 0, data);

    // draw code on surface
    uint8_t cell;

    uint32_t stride = cairo_format_stride_for_width(CAIRO_FORMAT_RGB24, qrcode.size);
    uint8_t* bmp = malloc(stride * qrcode.size);
    memset(bmp, 0, stride * qrcode.size);
    if (!bmp)
    {
        return NULL;
    }

    uint32_t* pixel = bmp;
    for (uint8_t y = 0; y < qrcode.size; ++y)
    {
        assert(y + 1 > y);
        for (uint8_t x = 0; x < qrcode.size; ++x)
        {
            assert(x + 1 > x);
            cell = 255 * (1 - (int)qrcode_getModule(&qrcode, x, y));
            
            pixel[x] = (0x0 << 24) | (cell << 16) | (cell << 8) | (cell << 0);
        }

        pixel = (uint32_t*)(bmp + stride * y);
    }

    surf = cairo_image_surface_create_for_data(
            bmp,
            CAIRO_FORMAT_RGB24,
            qrcode.size,
            qrcode.size,
            stride
    );
    
    if (cairo_surface_status(surf) != CAIRO_STATUS_SUCCESS)
    {
        printf("error: %s\n", cairo_status_to_string(cairo_surface_status(surf)));
        return NULL;
    }
    
    cairo_surface_flush(surf);
    cairo_surface_mark_dirty(surf);
    //cairo_surface_write_to_png(surf, "./out.png");

    return surf;
}

