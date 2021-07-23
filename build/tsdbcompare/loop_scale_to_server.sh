echo "run insert test"
for s in {100,200,400,600,800,1000};do
    for i in {2000,1000,500,1};do
		./write_to_server.sh -b $i -w 100 -g 1 -s $s
        for j in {50,16};do
            ./write_to_server.sh -b $i -w $j -g 0 -s $s
        done
    done
    for i in {2000,1000,500};do
    ./write_to_server.sh -b $i -w 1 -g 0 -s $s
    done
done
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 100
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 200
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 400
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 600
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 800
# ./write_to_server.sh -b 1 -w 1 -g 1 -s 1000
