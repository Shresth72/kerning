#include "utils.h"

#define NOB_IMPLEMENTATION
#define NOB_STRIP_PREFIX
#include "lib/nob.h"

#define WINDOW_WIDTH 600
#define WINDOW_HEIGHT 450
#define MAX_TABLES 64

void ParseFont(FILE *f, TableLoc *table_loc) {
  skip_bytes(f, 4);
  __uint16_t num_tables = read_uint16(f);
  assert(num_tables < MAX_TABLES);
  skip_bytes(f, 6);

  for (int i = 0; i < num_tables; ++i) {
    char tag[5];
    read_tag(f, tag);
    skip_bytes(f, 4);
    __uint32_t offset = read_uint32(f);
    skip_bytes(f, 4);

    strncpy(table_loc[i].tag, tag, 5);
    table_loc[i].offset = offset;
  }
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

  const char *fontPath = "../assets/JetBrainsMono-Bold.ttf";
  TableLoc table_loc[MAX_TABLES];

  FILE *fontFile = fopen(fontPath, "rb");
  if (!fontFile) {
    perror("Failed to open font file");
  }

  ParseFont(fontFile, table_loc);
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
