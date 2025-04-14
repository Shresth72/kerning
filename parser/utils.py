from typing import BinaryIO
import struct


def goto(file: BinaryIO, position: int) -> None:
    file.seek(position)


def skip_bytes(file: BinaryIO, num: int) -> None:
    file.seek(num, 1)


def read_uint8(file: BinaryIO) -> int:
    return struct.unpack(">b", file.read(1))[0]


def read_uint16(file: BinaryIO) -> int:
    return struct.unpack(">h", file.read(2))[0]


def read_uint32(file: BinaryIO) -> int:
    return struct.unpack(">i", file.read(4))[0]


def read_tag(file: BinaryIO) -> str:
    return file.read(4).decode("ascii")


def flag_bit_is_set(flag: int, bit_index: int):
    return (flag >> bit_index) & 1 == 1
