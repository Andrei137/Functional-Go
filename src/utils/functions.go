package functional

import (
    "math"
    "regexp"
    "strings"
    "strconv"
    "unicode"
)

// Helpers
func isVowel(ch rune) bool {
    return strings.ContainsRune("aăâeiîouAĂÂEIÎOU", ch)
}

func isConsonant(ch rune) bool {
    return !isVowel(ch)
}

func isLower(ch rune) bool {
    return ch >= 'a' &&
           ch <= 'z'
}

func isUpper(ch rune) bool {
    return ch >= 'A' &&
           ch <= 'Z'
}

func isDiacritic(ch rune) bool {
    return ch > unicode.MaxASCII
}

func decrypt(word string) string {
    var decrypted string = ""
    for _, ch := range word {
        decrypted += string('z' - ch + 'a')
    }
    return decrypted
}

func rymes(word, candidate string) bool {
    var (
        match int = 0
        i     int = len(word) - 1
        j     int = len(candidate) - 1
    )

    for i >= 0 && j >= 0 && word[i] == candidate[j] {
        match++
        i--
        j--
    }
    return match > 1
}

func isPerfectSquare(num int) bool {
    var sqrt float64 = math.Sqrt(float64(num))
    return sqrt == math.Floor(sqrt)
}

// Validators
func vowelsConsonants(word string) bool {
    var runes []rune = []rune(word)
    return CountSat(runes, isVowel) % 2 == 0 &&
           CountSat(runes, isConsonant) % 3 == 0
}

func isPalindrome(word string) bool {
    var left, right int = 0, len(word) - 1
    for left < right {
        if word[left] != word[right] {
            return false
        }
        left++
        right--
    }
    return true
}

func isGibberish(word string) bool {
    for i, ch := range word {
        if isVowel(ch) {
            if i < len(word) - 1 && word[i + 1] != 'p' {
                return false
            }
        }
    }
    return true
}

func vowelStartEnd(word string) bool {
    var runes []rune = []rune(word)
    return isVowel(runes[0]) &&
           isVowel(runes[len(word) - 1])
}

func anagramWithFacultate(word string) bool {
    ap := make(map[rune]int)
    for _, ch := range word {
        ap[ch]++
    }
    for _, ch := range "facultate" {
        if ap[ch] == 0 {
            return false
        }
        ap[ch]--
    }
    return true
}

func upperStartEnd(word string) bool {
    var runes []rune = []rune(word)
    return isUpper(runes[0]) &&
           isUpper(runes[len(word) - 1]) &&
           CountSat(runes, isLower) % 2 == 0
}

func containsDiacritics(word string) bool {
    return CountSat([]rune(word), isDiacritic) > 1
}

func hasSubstitutionPair(word string, words []string) bool {
    return AnySat(words, func(w string) bool {
        return decrypt(w) == word
    })
}

func hasRhymePair(word string, words[]string) bool {
    return AnySat(words, func(w string) bool {
        return word != w && rymes(word, w)
    })
}

func alternantVowelConsonant(word string) bool {
    if word == "" {
        return false
    }
    var runes []rune = []rune(word)
    if isVowel(runes[0]) {
        return false
    }
    for i := 0; i < len(runes) - 1; i++ {
        if isVowel(runes[i]) == isVowel(runes[i + 1]) {
            return false
        }
    }
    return true
}

func isStrongPassword(word string) bool {
    return regexp.MustCompile(`[a-z]`).MatchString(word) &&
           regexp.MustCompile(`[A-Z]`).MatchString(word) &&
           regexp.MustCompile(`\d`).MatchString(word) &&
           regexp.MustCompile(`[\W_]`).MatchString(word)

}

func isUnixPath(word string) bool {
    return strings.HasPrefix(word, "/") ||
           strings.HasPrefix(word, "./") ||
           strings.HasPrefix(word, "../")
}

func isRomanianName(word string) bool {
    return strings.HasSuffix(word, "escu")
}

func isFibonacciNumber(num string) bool {
    number, err := strconv.Atoi(num)
    if err != nil {
        return false
    }

    return isPerfectSquare(5 * number * number + 4) ||
           isPerfectSquare(5 * number * number - 4)
}

func hasThreeSetBits(num string) bool {
    number, err := strconv.Atoi(num)
    if err != nil {
        return false
    }

    var count int = 0
    for number > 0 {
        count += number & 1
        number >>= 1
    }
    return count == 3
}

// Main function
func GenericMapper(words []string, name string) int {
    validators := map[string]func(string) bool{
        "vowels_consonants"        : vowelsConsonants,
        "palindrome"               : isPalindrome,
        "gibberish"                : isGibberish,
        "vowel_start_end"          : vowelStartEnd,
        "anagram_with_facultate"   : anagramWithFacultate,
        "upper_start_end"          : upperStartEnd,
        "contains_diacritics"      : containsDiacritics,
        "has_substitution_pair"    : func(word string) bool {
            return hasSubstitutionPair(word, words)
        },
        "has_rhyme_pair"           : func(word string) bool {
            return hasRhymePair(word, words)
        },
        "alternant_vowel_consonant": alternantVowelConsonant,
        "strong_passwords"          : isStrongPassword,
        "unix_paths"                : isUnixPath,
        "romanian_names"            : isRomanianName,
        "fibonacci_numbers"         : isFibonacciNumber,
        "three_set_bits"            : hasThreeSetBits,
    }
    return ReduceSum(MapSat(words, validators[name]))
}