// Licensed under the GPL-v3
// Copyright: Sebastian Tilders <info@informatikonline.net> (c) 2019

#include <ncurses.h>
#include <stdlib.h>
#include <string.h>
#include <locale.h>

void del_window(WINDOW *win) {
    wbkgd(win,COLOR_PAIR(0));
    wclear(win);
    wrefresh(win);
    delwin(win);
}
void quit() {
    mvaddstr(0,0,"end");    
    refresh();
    endwin();
    exit(0);
}

void* bind_get_stdscr() {
    return stdscr;
}

void bind_set_locale() {
    setlocale(LC_ALL,"");
}
void bind_get_maxyx(WINDOW *win, int *y, int *x) {
    getmaxyx(win,*y,*x);
}

void bind_waddstr(WINDOW *win, char* str) {
    waddstr(win,str);
    free(str);
}

void bind_color_set(WINDOW* win, short pair) {
    wcolor_set(win,pair,0);
}

void bind_wbkgd(WINDOW *win, short pairId) {
    wbkgd(win,COLOR_PAIR(pairId));
}

int bind_wgetnstr(WINDOW *win, int max, char** str) {
    char* buf = calloc(1,max);
    wgetnstr(win,buf,max);
    *str = buf;
    return (int) strlen(buf);
}

int bind_fkey(int i) {
    return KEY_F(i);
}