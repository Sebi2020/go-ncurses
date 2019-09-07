// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

#include <curses.h>
#include <stdlib.h>
void bind_quit();
void bind_get_maxyx(WINDOW* win,int *y, int *x);
void bind_waddstr(WINDOW* win, char *str);
void* bind_get_stdscr();
void bind_set_locale();
void bind_color_set(WINDOW* win, short pair);
void bind_wbkgd(WINDOW *win, short pairId);
int bind_wgetnstr(WINDOW *win, int max, char** str);
int bind_fkey(int i);