package unittestsdemo

import "testing"

func Add(a, b int) int {
    return a + b
}

func TestAdd(t *testing.T) {
    // Table of test cases
    tests := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"positive numbers", 2, 2, 4},
        {"negative numbers", -1, -2, -3},
        {"zeroes", 0, 0, 0},
    }

    // Loop through the table and run tests
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.expected {
                t.Errorf("Unit Test %s Failed, Add(%d, %d) = %d; expected %d",tc.name, tc.a, tc.b, got, tc.expected)
            }
        })
    }
}