# output as png image
set terminal png

set output "std_5000_500.png"

# graph title
set title "ab -n 5000 -c 500 -g out.data http://localhost:8000/api/v1/workspaces/1/boards (STD)"

#nicer aspect ratio for image size
set size 0.95,1

# y-axis grid
set grid y

#x-axis label
set xlabel "request"

#y-axis label
set ylabel "response time (ms)"

#plot data from "out.data" using column 9 with smooth sbezier lines
plot "out.data" using 9 smooth sbezier with lines title "MyTrello"
