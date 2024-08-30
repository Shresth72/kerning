import matplotlib.pyplot as plt

points = [
(1, 64836),
(44031, 64836),
(44031, 65516),
(65819, 65516),
(66004, 130743),
(65818, 131052),
(65645, 196258),
(65826, 261491),
(65826, 262096),
(65620, 327330),
(111893, 327630),
(111893, 392565),
(112078, 392845),
(112077, 458073),
(141263, 458073),
(141451, 458073),
]

contour_end_indices = [3, 6, 9, 12, 15]

plt.figure(figsize=(10, 10))
plt.title('Contour Plot')
plt.xlabel('X')
plt.ylabel('Y')

contour_start_index = 0
for contour_end_index in contour_end_indices:
    contour_points = points[contour_start_index:contour_end_index + 1]
    x_values, y_values = zip(*contour_points)
    
    plt.plot(x_values, y_values, marker='o', color='red')
    
    contour_start_index = contour_end_index + 1

plt.savefig('contour_plot.png')

plt.show()
