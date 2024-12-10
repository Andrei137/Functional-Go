package functional

import (
    "math"
    "sync"
    "go.uber.org/zap"
    "github.com/fatih/color"
)

type TestCase struct {
    Name     string
    Queries  [][]string
    Expected float64
}

func almostEqual(a, b float64) bool {
    const epsilon float64 = 1e-2
    return math.Abs(a - b) <= epsilon
}

func mapReduce(data [][]string, testName string) float64 {
    var (
        wg     sync.WaitGroup
        length int            = len(data)
        ch     chan int       = make(chan int, length)
    )

    for _, words := range data {
        wg.Add(1)
        go func(words []string) {
            defer wg.Done()
            ch <- GenericMapper(words, testName)
        }(words)
    }
    go func() {
        wg.Wait()
        close(ch)
    }()

    return float64(ReduceSum(ch)) / float64(length)
}

func runTest(logger *zap.SugaredLogger, testIndex int, testName string, data [][]string) float64 {
    logger.Infof(color.HiBlueString("Test cerinta %d (%s)", testIndex, testName))
    logger.Infof("Input %v", data)
    var result float64 = mapReduce(data, testName)
    logger.Infof("Output %.3f\n", result)
    return result
}

func RunTests() {
    var (
        logger     *zap.SugaredLogger = InitLogger().Sugar()
        failed     []int
    )
    defer logger.Sync()

    tests := []TestCase {
        {
            Name   : "vowels_consonants",
            Queries: [][]string{
                {"aabbb", "ebep", "blablablaa", "hijk", "wsww"},
                {"abba", "eeeppp", "cocor", "ppppppaa", "qwerty", "acasq"},
                {"lalala", "lalal", "papapa", "papap"},
            },
            Expected: 2.33,
            // eeeppp nu e valid, dar blablablaa este
        },
        {
            Name   : "palindrome",
            Queries: [][]string {
                {"a1551a", "parc", "ana", "minim", "1pcl3"},
                {"calabalac", "tivit", "leu", "zece10", "ploaie", "9ana9"},
                {"lalalal", "tema", "papa", "ger"},
            },
            Expected: 2.33,
        },
        {
            Name   : "gibberish",
            Queries: [][]string{
                {"apap", "paprc", "apnap", "mipnipm", "copil"},
                {"cepr", "program", "lepu", "zepcep", "golang", "tema"},
                {"par", "impar", "papap", "gepr"},
                // se pune lepu?
            },
            Expected: 3.0,
        },
        {
            Name   : "vowel_start_end",
            Queries: [][]string{
                {"ana", "parc", "impare", "era", "copil"},
                {"cer", "program", "leu", "alee", "golang", "info"},
                {"inima", "impar", "apa", "eleve"},
            },
            Expected: 2.66,
        },
        {
            Name   : "anagram_with_facultate",
            Queries: [][]string{
                {"acultatef", "parc", "cultateaf", "faculatet", "copil"},
                {"cer", "tatefacul", "leu", "alee", "golang", "ultatefac"},
                {"tefaculta", "impar", "apa", "eleve"},
            },
            Expected: 2.00,
        },
        {
            Name   : "upper_start_end",
            Queries: [][]string{
                {"AcasA", "CasA", "FacultatE", "SisTemE", "distribuite"},
                {"GolanG", "map", "reduce", "Problema", "TemA", "ProieCt"},
                {"LicentA", "semestru", "ALGORitM", "StuDent"},
            },
            Expected: 1.66,
        },
        {
            Name   : "contains_diacritics",
            Queries: [][]string{
                {"țânțar", "carte", "ulcior", "copac", "plante"},
                {"beci", "", "mlăștinos", "astronaut", "stele", "planete"},
                {"floare", "somn", "șosetă", "scârțar"},
            },
            Expected: 1.33,
        },
        {
            Name   : "has_substitution_pair",
            Queries: [][]string{
                {"prajitura", "camion", "foaie", "liliac", "ulzrv"},
                {"carte", "trofeu", "xzigv", "laptop", "scris", "muzica"},
                {"pictura", "telefon", "parapanta", "catel"},
                // foaie ~ ulzrv (nu uezrv), si carte ~ xzigv (nu xzigw)
            },
            Expected: 1.33,
        },
        {
            Name   : "has_rhyme_pair",
            Queries: [][]string{
                {"stele", "mele", "borcan", "vajnic", "strașnic"},
                {"crocodil", "garnisit", "mușețel", "făurit", "arhanghel", "noapte"},
                {"lampă", "sine", "cine", "toriște"},
                // mușețel ar trebui sa rimeze cu arhanghel
                // daca garnisit rimeaza cu făurit
            },
            Expected: 3.33,
        },
        {
            Name   : "alternant_vowel_consonant",
            Queries: [][]string{
                {"caracatiță", "ceva", "saar", "aaastrfb", ""},
                {"aaabbbccc", "caporal", "ddanube", "jahfjksgfjhs", "ajsdas", "urs"},
                {"scoica", "coral", "arac", "karnak"},
            },
            Expected: 1.33,
        },
        {
            Name   : "strong_passwords",
            Queries: [][]string{
                {"sadsa1@A", "cevaA!4", "saar", "aaastrfb", ""},
                {"aaabbbccc", "!Caporal1", "ddanube", "jahfjksgfjhs", "ajsdas","urs"},
                {"scoica", "Coral!@12", "arac", "karnak"},
            },
            Expected: 1.33,
        },
        {
            Name   : "unix_paths",
            Queries: [][]string{
                {"/dev/null", "/bin", "saar", "teme/scoala/2020", ""},
                {"proiect/tema", "/dev", "ddanube", "jahfjksgfjhs", "ajsdas","urs"},
                {"scoica", "/teme/repos/git", "arac", "karnak"},
            },
            Expected: 1.33,
        },
        {
            Name   : "romanian_names",
            Queries: [][]string{
                {"Popescu", "Ionescu", "Pop", "aaastrfb", ""},
                {"Nicolae", "Dumitrescu", "ddanube", "jahfjksgfjhs", "ajsdas","urs"},
                {"Dumitru", "Angelescu", "arac", "karnak"},
            },
            Expected: 1.33,
        },
        {
            Name   : "fibonacci_numbers",
            Queries: [][]string{
                {"12", "5", "6", "13", "7"},
                {"21", "20", "42", "43", "8","38"},
                {"54", "55", "34", "100"},
            },
            Expected: 2.00,
            // 55 is a fibonacci number
        },
        {
            Name   : "three_set_bits",
            Queries: [][]string{
                {"1", "13", "6", "7", "9"},
                {"19", "20", "43", "43", "21","53"},
                {"54", "55", "28", "101"},
            },
            Expected: 1.66,
        },
    }

    for i, test := range tests {
        if !almostEqual(runTest(logger, i + 1, test.Name, test.Queries), test.Expected) {
            failed = append(failed, i + 1)
        }
    }
    if len(failed) == 0 {
        logger.Infof(color.HiGreenString("All tests passed!"))
    } else {
        logger.Errorf(color.HiRedString("Failed tests: %v", failed))
    }
}
