set -u

color_ok="\e[32;1m"
color_failed="\e[31;1m"
color_off="\e[m"

tmp_dir=/tmp/tmpfs/out
asm_dir=$tmp_dir/asm
bin_dir=$tmp_dir/bin
mkdir -p $tmp_dir $asm_dir $bin_dir

function run_unit_test {
    test_name=$1
    src_file=tests/$test_name/main.go
    asm_file=$asm_dir/$test_name.s
    bin_file=$bin_dir/$test_name

    make indigo --silent
    ./output/indigo tests/$test_name/main.go > $asm_file
    clang -o $bin_file $asm_file

    expected_file=tests/$test_name/expected.txt
    actual_file=$tmp_dir/actual.txt
    $bin_file
    echo $? > $actual_file

    if cmp -s $actual_file $expected_file; then
        echo "[ok]     ${test_name}"
        return 0
    else
        echo "[failed] ${test_name}"
        return 1
    fi
}

function run_all_test {
    for test in ./tests/*; do
        test_name=$(basename $test)
        run_unit_test $test_name
    done
}

run_all_test
