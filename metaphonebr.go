package metaphonebr

import (
//   "fmt"
   "errors"
   "regexp"
   "strings"
   "github.com/arbovm/levenshtein"
   "github.com/luisfurquim/goose"
)

type NameT struct {
   Words []string
   Mtfs []string
}

type ruleT struct {
   re *regexp.Regexp
   tgt string
   offset int
   inc    int
   isFirst bool
}

var rules []ruleT
var reVowell *regexp.Regexp
var repAccent *strings.Replacer
var reWord *regexp.Regexp

var Verbose bool = true

var LevThreshold float32 = .5 // Margem percentual para considerar dois metaphones (MTFs) similares usando levenshtein - 1 (min == 1). Vide func IsSim(...)

var Goose goose.Alert

var preps map[string]bool = map[string]bool{
   "DE": true,
   "DO": true,
   "DA": true,
   "DOS": true,
   "DAS": true,
}

var ErrInvalidName = errors.New("Invalid name")


func (n NameT) String() string {
   var s string

   s = strings.Join(n.Words," ")

   if Verbose {
      s += " (" + strings.Join(n.Mtfs," ") + ")"
   }

   return s
}

func init() {
   reWord = regexp.MustCompile("(\\pL+)")
   // SOURCE: http://sourceforge.net/p/metaphoneptbr/code/ci/master/tree/README
   rules = []ruleT{
      ruleT{
         regexp.MustCompile("(?i)^a"),
         "A",
         0,
         1,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^[ei]"),
         "I",
         0,
         1,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^[ou]"),
         "U",
         0,
         1,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^b"),
         "B",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^c(?:(?:[bcdfgjklmnpqrstvwxzaou])|$)"),
         "K",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^chr"),
         "KR",
         0,
         3,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^c[ei]"),
         "S",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^d"),
         "D",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^f"),
         "F",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^g[aou]"),
         "G",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^gh[bcdfgjklmnpqrstvwxz]"),
         "G",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^g[ei]"),
         "J",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^gh[ei]"),
         "J",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ha"),
         "A",
         0,
         2,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^h[ei]"),
         "I",
         0,
         2,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^h[ou]"),
         "U",
         0,
         2,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^lh"),
         "1",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^nh"),
         "3",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^h"),
         "",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^j"),
         "J",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^k"),
         "K",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^l[aou]"),
         "l",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^m"),
         "M",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^n$"),
         "M",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ph"),
         "F",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^p"),
         "P",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^q"),
         "K",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^qu"),
         "K",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^r"),
         "2",
         0,
         1,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^r$"),
         "R",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^rr"),
         "2",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[aou]r[aeiou]"),
         "R",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^.r[bcdfghjklmnpqrstvwxz]"),
         "R",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[bcdfghjklmnpqrstvwxz]r[aeiou]"),
         "R",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ss"),
         "S",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^(?:s|c)h"),
         "X",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^sch"),
         "X",
         0,
         3,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^sc[ei]"),
         "S",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^sc"),
         "SK",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^s[bdfgjklmnpqrstvwxz]"),
         "S",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^t"),
         "T",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^th"),
         "T",
         0,
         2,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^v"),
         "V",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^w[lraeiou]"),
         "V",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^w[bcdfghjklmnpqrstvwxz]"),
         "",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^x$"), // REVISAR!!!
         "X",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ex[aeiou]"),
         "Z",
         -1,
         1,
         true,
      },

      ruleT{
         regexp.MustCompile("(?i)^ex[ei]"),
         "X",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ex[ptc]"),
         "S",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^.ex[aou]"),
         "X",
         -2,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ex[aou]"),
         "KS",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^ex."),
         "KS",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[aeiouckglrx][aiou]x"),
         "X",
         -2,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[dfmnpqstvz][aou]x"),
         "KS",
         -2,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[aeiou][i][aeiou]"),
         "I",
         -1,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^[y]"),
         "I",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^Z$"),
         "S",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^Z"),
         "Z",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^X"),
         "X",
         0,
         1,
         false,
      },

      ruleT{
         regexp.MustCompile("(?i)^S"),
         "S",
         0,
         1,
         false,
      },

   }

   reVowell = regexp.MustCompile("(?i)[aeiou]")
   repAccent = strings.NewReplacer(
      "ç", "ss",
      "Ç", "ss",
      "á", "a",
      "é", "e",
      "í", "i",
      "ó", "o",
      "ú", "u",
      "Á", "a",
      "É", "e",
      "Í", "i",
      "Ó", "o",
      "Ú", "u",
      "ã", "a",
      "ẽ", "e",
      "ĩ", "i",
      "õ", "o",
      "ũ", "u",
      "Ã", "a",
      "Ẽ", "e",
      "Ĩ", "i",
      "Õ", "o",
      "Ũ", "u",
      "â", "a",
      "ê", "e",
      "î", "i",
      "ô", "o",
      "û", "u",
      "Â", "a",
      "Ê", "e",
      "Î", "i",
      "Ô", "o",
      "Û", "u",
      "à", "a",
      "è", "e",
      "ì", "i",
      "ò", "o",
      "ù", "u",
      "À", "a",
      "È", "e",
      "Ì", "i",
      "Ò", "o",
      "Ù", "u",
      "ä", "a",
      "ë", "e",
      "ï", "i",
      "ö", "o",
      "ü", "u",
      "Ä", "a",
      "Ë", "e",
      "Ï", "i",
      "Ö", "o",
      "Ü", "u",
      "ý", "y",
      "ỳ", "y",
      "ỹ", "y",
      "ŷ", "y",
      "ÿ", "y",
      "ñ", "n")
}

func Pack(s string) string {
   var ret string
   var inc int

   s = repAccent.Replace(s)

   for i:=0; i<len(s); {
//      for j,rule := range rules {
      for _,rule := range rules {
         if rule.isFirst && (i>0) {
            continue
         }
         if (i+rule.offset) < 0 {
            continue
         }
         if rule.re.MatchString(s[i+rule.offset:]) {
            ret += rule.tgt
            inc  = rule.inc
//            fmt.Printf("Rule %d %#v, input:%s %#v",j,rule,string(b[i+rule.offset:]),b[i+rule.offset:])
            break
         }
      }
      if inc>0 {
         i += inc
      } else {
//         fmt.Printf("No rule found for %c %x\n",b[i],b[i])
         i++
      }
   }

   if ret == "" {
      ret = strings.ToUpper(reVowell.ReplaceAllString(s,""))
//      fmt.Printf("%s %#v     %s\n",ret,[]byte(ret),s)
   }

   return ret
}

func Parse(nm string) *NameT {
   var i int
   var w []string
   var ret NameT

   ret = NameT{}
   words := reWord.FindAllStringSubmatch(strings.ToUpper(nm),-1)
   if words == nil {
      return nil
   }

   ret.Words = make([]string,len(words))
   ret.Mtfs = make([]string,len(words))
   for i, w = range words {
      ret.Words[i] = w[1]
      ret.Mtfs[i]  = Pack(w[1])
   }
//   fmt.Printf("%#v",ret)
   return &ret
}


func IsSim(mtf1, mtf2 string) bool {
   var margem int

   // margem = min(|mtf1|, |mtf2|)
   margem = len(mtf1)
   if margem > len(mtf2) {
      margem = len(mtf2)
   }

   margem-- // margem > 1 somente se ambos os MTFs tiverem pelo menos 5 caracteres...
   margem = int(LevThreshold * float32(margem))
   if margem < 1 {
      margem = 1
   }

   if levenshtein.Distance(mtf1, mtf2) <= margem {
      return true
   }
   return false
}


func WordSim(w1, w2 string) float32 {
   var maxsz float32

   if len(w1) < len(w2) {
      maxsz = float32(len(w2))
   } else {
      maxsz = float32(len(w1))
   }

   return 1.0 - (float32(levenshtein.Distance(w1,w2)) / maxsz)
}

func (n NameT) SimString(name string) (float32, error) {
   var pes2              *NameT

   pes2 = Parse(name)
   if pes2 == nil {
      return -1, ErrInvalidName
   }

   return n.Sim(pes2)
}

func (n NameT) Sim(pes2 *NameT) (float32, error) {
   var i, j, pos          int
   var mtf1, mtf2         string
   var sim, simtotal      float32
   var matches            int
   var qtPrep, qtPrep2    int
//   var err                error

   simtotal = 1.0

   for i, mtf1 = range n.Mtfs {
      if qtPrep2>0 {
         Goose.Logf(4,"DROP qtPrep2: %d\n",qtPrep2)
         qtPrep2 = 0
      }
      if preps[n.Words[i]] {
         qtPrep++
         Goose.Logf(4,"qtPrep: %d",qtPrep)
         continue
      }
      sim = 0.0

      for j = pos; j < len(pes2.Mtfs); j++ {
         if preps[pes2.Words[j]] {
            qtPrep2++
            continue
         }
         mtf2 = pes2.Mtfs[j]
         if mtf1==mtf2 {
            sim = WordSim(n.Words[i],pes2.Words[j])
            break
         } else if IsSim(mtf1, mtf2) {
            sim = WordSim(n.Words[i],pes2.Words[j])
            break
         }
      }
      if sim > 0.0 {
         Goose.Logf(4,"Match %f, i:%d, j:%d, pos=%d, newpos=%d",sim,i,j,pos,j+1)
         simtotal *= sim
//         fmt.Printf("Sim: %f, Total: %f\n",sim,simtotal)
         matches++
         pos = j + 1
         if qtPrep2>0 {
            Goose.Logf(4,"Adding qtprep2: %d",qtPrep2)
            qtPrep += qtPrep2
            Goose.Logf(4,"ACC qtPrep2: %d\n",qtPrep2)
         }
      }
   }
   Goose.Logf(4,"Matches: %d, len1:%d, len2:%d\n",matches,len(n.Words),len(pes2.Words))
   Goose.Logf(4,"qtPrep: %d\n",qtPrep)
   sim = simtotal * (2.0*float32(matches)/float32(len(n.Words)+len(pes2.Words)-qtPrep))
   Goose.Logf(4,"%f: %s x %s\n",sim,n,pes2)
   return sim, nil
}


