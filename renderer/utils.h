#pragma once
#include <errno.h>
#include <math.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>

#include "lib/nob.h"
#include "lib/stb_ds.h"

#include "lib/raylib/raylib-5.5_linux_amd64/include/raylib.h"
#include "lib/raylib/raylib-5.5_linux_amd64/include/rlgl.h"

typedef struct {
  char tag[5];
  __uint32_t offset;
} GlyphTableEntry;

typedef struct {
  GlyphTableEntry *tables;
  int count;
} GlyphTableMap;

// GLYPH
void GetAllGlyphLocations(FILE *f, GlyphTableMap *glyph_map);

// UTILS
long get_loc(FILE *file);
void goto_loc(FILE *file, long position);
void skip_bytes(FILE *file, long num);

__uint8_t read_uint8(FILE *file);
__int16_t read_int16(FILE *file);
__uint16_t read_uint16(FILE *file);
__uint32_t read_uint32(FILE *file);
void read_tag(FILE *file, char *tag);

bool flag_bit_is_set(__uint32_t flag, int bit_index);
__uint32_t get_table_offset(GlyphTableMap *glyph_map, const char *tag);
