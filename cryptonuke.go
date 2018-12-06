/*
 * Czz78 for Ibic83
 * This is ICBM Atomic Bomb Crypt Fallout76 decryptor written in Golang for Fun
 * I have to learn Go so this is a real basic example
 */

package main

import (
    "fmt"
    "strings"
    "unicode"
    "io/ioutil"
    "os"
    "encoding/json"
    "flag"
    "regexp"
    "sort"
    "strconv"
)


/*
 *  Constants
 *  @ALPHABET The alphabet
 *  @WORDLIST Name of word list file
 *  @WORDLIST Name of key word list file
 */
const ALPHABET string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const WORDLIST string = "wordlist.json"
const KEYWORDLIST string = "keywordlist.json"
const LISTPATH string = "lists"

/*
 *  Code Pair
 *  @letters string Key letters
 *  @numbers []int Key numbers
 */
type Code struct {
    letters string
    numbers []int
}


/*
 *  Cipher Keyword
 *  @cipherword string Cipher word
 *  @cipherstring string Cipher string
 *  @ciperwordlen int Lenght of cipher word
 *  @deccyptedKey Code The decripted code pair
 */
type CipherKeyWord struct {
     cipherword string
     cipherstring string
     cipherwordlen int
     decryptedKey Code
}

/*
 *  Solution struct
 *  @password []int Password for Nuke
 *  @keyword string The Key word for Nuke
 *  @codeword string The code word for Nuke
 */
type Solution struct {
     password []int
     keyword string
     codeWord string
}


/*
 * Decrypts a nuclear key
 * @key string Private key
 * @code Code  Key pair code
 * @return Code Key pair decrypted code
 */
func decryptKey(key string, code Code) Code {

  var let string
  var num []int

  for i := 0; i < len(code.letters); i++ {

      pos := strings.Index(key, string(code.letters[i]))
      let = let + string(ALPHABET[pos])
      num = append(num ,int(code.numbers[i]))

  }

  return Code {
               letters : let,
               numbers : num,
              }

}


/*
 * Cipher string from a word
 * @word string The word
 * @return string Cipher
 */
func returnCipherString (word string) string {

    // we could check if string passed is only composed of alphabet upper case  letters, but we don't

    var buffer string = strings.ToUpper(strings.Map(func(r rune) rune {
                            if unicode.IsSpace(r) {
                                return -1
                            }
                            return r
                        }, word))

    var char string
    for i := 0; i < len(ALPHABET) ; i++ {
        char = string(ALPHABET[i])
        if !strings.ContainsAny(buffer, char) {
            buffer = buffer + char
        }
    }

    return buffer;
}


/*
 * Find Cipher Key Words From a list
 * @word []string Slice of keywords
 * @code Code Key pair code
 * @return []CipherKeyWord Struct Slices of cipher key words
 */
func cipherKeyWords (words []string, code Code ) []CipherKeyWord {

    result := []CipherKeyWord{}
    var privateKey string
    var d Code

    for i:=0; i< len(words); i++ {

        privateKey = returnCipherString(words[i])
        d = decryptKey(privateKey, code)

        ckw := CipherKeyWord {
                             cipherword : words[i],
                             cipherstring : privateKey,
                             cipherwordlen : len(words[i]),
                             decryptedKey : d,
                            }

        result = append(result, ckw)

    }

    return result

}


/*
 *  Get a list from a json file
 *  @p string path of the file
 *  @return []string
 */
func getJsonList(p string) []string {

    var res []string
    file, err := ioutil.ReadFile(p)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    err = json.Unmarshal([]byte(file), &res)
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }

    return res

}


/*
 * Filter a key word list
 * @s string Regex passed ( only a-z and . are allowed for this program )
 * @list []string A list of words
 */
func getKeywordList(s string, list []string) []string {

    var res []string 

    for i := 0 ; i < len(list); i++ {

       matched, _ := regexp.MatchString(s, list[i])
       if matched {
           res = append(res, list[i])
       }

    }
    return res

}


/*
 * Sort cipher only works with golang 1.7
 */
type ByRune []rune
func (r ByRune) Len() int           { return len(r) }
func (r ByRune) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRune) Less(i, j int) bool { return r[i] < r[j] }
func sortCipher(c string) string {

    var r []rune
      for _, runeValue := range c {
              r = append(r, runeValue)
      }
    var p ByRune = r
    sort.Sort(p)



    return string(p)

}


/*
 * Sort cipher only works on golang 1.8
 */
/*
func sortCipher(c string) string {

    var r []rune
    for _, runeValue := range c {
        r = append(r, runeValue)
    }

    sort.Slice(r, func(i, j int) bool {
              return r[i] < r[j]
    })

    return string(r)

}
*/

/*
 * Compares to ciphers if they have same letters
 * @c1 string Cipher 1
 * @c2 string Cipher 2
 * @return bool True or False
 */
func compareCiphers(c1 string, c2 string) bool{

    if sortCipher(c1) == sortCipher(c2) {
        return true
    }
    return false

}


/*
 * Get The Solution
 * @cipher []CipherKeyWord Slice of cipher keywords
 * @list []string A list of ciphers
 * @return Solution The solution: password, keyword, codeword
 */
func getCodeFromCiphers(cipher []CipherKeyWord, list []string) []Solution {

    result := []Solution{}

    for i := 0; i < len(cipher); i++ {

        c := cipher[i].decryptedKey

        for j:=0; j< len(list); j++ {

            l := list[j]

            if len(l) == 8 {

                if compareCiphers(l, c.letters) {

                    var code []int
                    for k :=0; k<len(l); k++ {
                        pos := strings.Index(c.letters,string(l[k]))
                        code = append(code,c.numbers[pos])
                    }

                    ckw := Solution {
                                     password : code,
                                     keyword : l, 
                                     codeWord: cipher[i].cipherword,
                                    }
                    result = append(result, ckw)

                }
            }
        }
    }

    return result

}


/*
 *  The main function
 */
func main() {


    wordlist :=  getJsonList(LISTPATH + "/" + WORDLIST)
    keywordlist := getJsonList(LISTPATH + "/" + KEYWORDLIST)


    /*
     * Flags
     */
    var knowntext string
    var letters string
    var numbers string
    flag.StringVar(&knowntext, "knowntext", "", "Insert the Known Text with dots where you don't have a letter es: B...ST")
    flag.StringVar(&letters, "letters", "", "Insert the key letters es: DEFMPRST")
    flag.StringVar(&numbers, "numbers", "", "Insert the key numbers es: 48110398")

    flag.Parse()

    if knowntext == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    if letters == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    if numbers == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    var num []int
    for i :=0; i< len(numbers); i++ {
         n, err := strconv.Atoi(string(numbers[i]))
         if err != nil {
             fmt.Println("Error inserting numbers")
             flag.PrintDefaults()
             os.Exit(2)
         }
         num = append(num, n)


    }

    code := Code  {letters: letters , numbers: num }

    keywordlist = getKeywordList(knowntext, keywordlist)
    s:=cipherKeyWords(keywordlist,code)

    fmt.Println(getCodeFromCiphers(s, wordlist))

}





