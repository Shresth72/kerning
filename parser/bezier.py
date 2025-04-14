class Point:
    def __init__(self, x: int, y: int, on_curve: bool):
        self.x = x
        self.y = y
        self.on_curve = on_curve

    def __repr__(self):
        return f"({self.x}, {self.y})"


def DrawLine(p1: Point, p2: Point, ax):
    ax.plot([p1.x, p2.x], [p1.y, p2.y], "k-")


def DrawPoint(p: Point, ax):
    ax.plot(p.x, p.y, "ro", markersize=3)


def LinearInterpolation(start: Point, end: Point, time: float) -> Point:
    return Point(
        x=start.x + (end.x - start.x) * time,
        y=start.y + (end.y - start.y) * time,
        on_curve=False,
    )


def BezierInterpolation(p0: Point, p1: Point, p2: Point, time: float) -> Point:
    intermediateA = LinearInterpolation(p0, p1, time)
    intermediateB = LinearInterpolation(p1, p2, time)
    return LinearInterpolation(intermediateA, intermediateB, time)


def DrawBezier(p0: Point, p1: Point, p2: Point, resolution: int, ax) -> None:
    prev_point_on_curve = p0

    for i in range(resolution):
        t = (i + 1.0) / resolution
        next_point_on_curve = BezierInterpolation(p0, p1, p2, t)
        DrawLine(prev_point_on_curve, next_point_on_curve, ax)
        prev_point_on_curve = next_point_on_curve
