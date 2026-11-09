fn linear_search(arr: &[i64], value: i64) -> i64 {
    for (i, val) in arr.iter().enumerate() {
        if *val == value {
            return i64::try_from(i).expect("i convert failed");
        }
    }
    return -1;
}

#[test]
fn test_linear_search() {
    assert_eq!(linear_search(&[10, 11, 12], 10), 0);
    assert_eq!(linear_search(&[10, 11, 12], 13), -1);
    assert_eq!(linear_search(&[10, 11, 12], 11), 1);
    assert_eq!(linear_search(&[10, 11, 12], 12), 2);
    assert_eq!(linear_search(&[12, 11, 10], 12), 0);
    assert_eq!(linear_search(&[11, 12, 12, 12], 12), 1);
}
