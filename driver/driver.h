#ifndef DRIVER_H
#define DRIVER_H

#ifdef __cplusplus
extern "C"
{
#endif

int driver_create_device(void);
int driver_mouse_rel(int, int);
int driver_mouse_btn(int, int, int);
int driver_report(void);
void driver_destroy_device(void);

#ifdef __cplusplus
}
#endif

#endif /* DRIVER_H */

