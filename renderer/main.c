#include "utils.h"

#define NOB_IMPLEMENTATION
#define NOB_STRIP_PREFIX
#include "lib/nob.h"

#define STB_DS_IMPLEMENTATION
#include "lib/stb_ds.h"

#define WINDOW_WIDTH 1200
#define WINDOW_HEIGHT 400
#define MAX_TABLES 64

void ParseFont(FILE *f) {
  skip_bytes(f, 4);
  __uint16_t num_tables = read_uint16(f);
  skip_bytes(f, 6);

  TagOffsetEntry *tables = malloc(sizeof(TagOffsetEntry) * num_tables);
  for (int i = 0; i < num_tables; ++i) {
    char tag[5];
    read_tag(f, tag);
    skip_bytes(f, 4);
    __uint32_t offset = read_uint32(f);
    skip_bytes(f, 4);

    strncpy(tables[i].tag, tag, 5);
    tables[i].offset = offset;
  }

  TagOffsetMap tag_offset_map = {.tables = tables, .count = num_tables};
  __uint32_t *all_glyph_locs = GetAllGlyphLocations(f, &tag_offset_map);

  GlyphUnicodeIndexMap *unicode_index_map =
      GetUnicodeGlyphIndexMap(f, &tag_offset_map);
  printf("Length of unicode_index_map: %d", unicode_index_map->count);

  free(tables);
  free(all_glyph_locs);
  free(unicode_index_map->indices);
  free(unicode_index_map);
}

int main() {
  const char *glyphs_dir = "glyfs";
  if (mkdir(glyphs_dir, 0755) == -1) {
    if (errno == EEXIST) {
    } else {
      perror("mkdir");
      return 1;
    }
  } else {
    printf("Directory created: %s\n", glyphs_dir);
  }

  const char *fontPath = "../assets/Meditative.ttf";
  FILE *fontFile = fopen(fontPath, "rb");
  if (!fontFile) {
    perror("Failed to open font file");
  }

  ParseFont(fontFile);
  fclose(fontFile);

  // InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "RandomArt");
  // SetTargetFPS(60);
  // while (!WindowShouldClose()) {
  //   BeginDrawing();
  //
  //   DrawRectangle(WINDOW_WIDTH / 2, WINDOW_HEIGHT / 2, 100, 100, GREEN);
  //
  //   EndDrawing();
  // }

  return 0;
}
