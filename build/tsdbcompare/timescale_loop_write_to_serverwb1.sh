for s in {100,200,400,600,800,1000};do
  ./timescale_write_to_server.sh -b 1 -w 1 -g 1 -s $s
  for i in {5000,2500,1000};do
    ./timescale_write_to_server.sh -b $i -w 1 -g 0 -s $s
    for j in {16,50,100};do
      ./timescale_write_to_server.sh -b 1 -w $j -g 0 -s $s
        done
    done
done
