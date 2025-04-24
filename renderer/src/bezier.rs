pub type Point = [f64; 2];

pub fn draw_bezier(p0: &Point, p1: &Point, p2: &Point, res: i32) {
    let mut prev = bezier_interpolation(p0, p1, p2, 0.0);

    for i in 1..=res {
        let t = (i as f64 + 1.0) / res as f64;
        let next = bezier_interpolation(p0, p1, p2, t);
        draw_line(&prev, &next);
        prev = next;
    }
}

fn bezier_interpolation(p0: &Point, p1: &Point, p2: &Point, t: f64) -> Point {
    let a = lerp(p0, p1, t);
    let b = lerp(p1, p2, t);
    lerp(&a, &b, t)
}

fn lerp(a: &Point, b: &Point, t: f64) -> Point {
    [a[0] + (b[0] - a[0]) * t, a[1] + (b[1] - a[1]) * t]
}

fn draw_line(p1: &Point, p2: &Point) {
    println!("line bw: {:?}, {:?}", p1, p2);
}

fn draw_point(p: &Point) {
    println!("point: {:?}", p);
}
