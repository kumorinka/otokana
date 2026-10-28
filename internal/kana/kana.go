// Package kana: tabel kana (PRD bagian 5) dan alias.
//
// Go murni, tanpa import mygo.
package kana

// Script is a writing system.
type Script int

const (
	Hiragana Script = iota
	Katakana
)

func (s Script) String() string {
	if s == Katakana {
		return "katakana"
	}
	return "hiragana"
}

// Group is a group of the table, in the order decks are built from.
type Group int

const (
	Mono  Group = iota // monographs
	Dia                // diacritics
	Di                 // digraphs
	DiDia              // digraphs + diacritics
)

// Entry is one kana with its accepted romaji.
type Entry struct {
	Romaji   string // main romaji, unique in the table, shown as the hint
	Hiragana string
	Katakana string
	Aliases  []string // accepted answers, Romaji included
	Group    Group
}

// Kana returns the entry in script.
func (e *Entry) Kana(s Script) string {
	if s == Katakana {
		return e.Katakana
	}
	return e.Hiragana
}

// Table holds the 107 entries, by group, each in gojūon order.
var Table = []Entry{
	{"a", "あ", "ア", []string{"a"}, Mono},
	{"i", "い", "イ", []string{"i"}, Mono},
	{"u", "う", "ウ", []string{"u"}, Mono},
	{"e", "え", "エ", []string{"e"}, Mono},
	{"o", "お", "オ", []string{"o"}, Mono},
	{"ka", "か", "カ", []string{"ka"}, Mono},
	{"ki", "き", "キ", []string{"ki"}, Mono},
	{"ku", "く", "ク", []string{"ku"}, Mono},
	{"ke", "け", "ケ", []string{"ke"}, Mono},
	{"ko", "こ", "コ", []string{"ko"}, Mono},
	{"sa", "さ", "サ", []string{"sa"}, Mono},
	{"shi", "し", "シ", []string{"shi", "si"}, Mono},
	{"su", "す", "ス", []string{"su"}, Mono},
	{"se", "せ", "セ", []string{"se"}, Mono},
	{"so", "そ", "ソ", []string{"so"}, Mono},
	{"ta", "た", "タ", []string{"ta"}, Mono},
	{"chi", "ち", "チ", []string{"chi", "ti"}, Mono},
	{"tsu", "つ", "ツ", []string{"tsu", "tu"}, Mono},
	{"te", "て", "テ", []string{"te"}, Mono},
	{"to", "と", "ト", []string{"to"}, Mono},
	{"na", "な", "ナ", []string{"na"}, Mono},
	{"ni", "に", "ニ", []string{"ni"}, Mono},
	{"nu", "ぬ", "ヌ", []string{"nu"}, Mono},
	{"ne", "ね", "ネ", []string{"ne"}, Mono},
	{"no", "の", "ノ", []string{"no"}, Mono},
	{"ha", "は", "ハ", []string{"ha"}, Mono},
	{"hi", "ひ", "ヒ", []string{"hi"}, Mono},
	{"fu", "ふ", "フ", []string{"fu", "hu"}, Mono},
	{"he", "へ", "ヘ", []string{"he"}, Mono},
	{"ho", "ほ", "ホ", []string{"ho"}, Mono},
	{"ma", "ま", "マ", []string{"ma"}, Mono},
	{"mi", "み", "ミ", []string{"mi"}, Mono},
	{"mu", "む", "ム", []string{"mu"}, Mono},
	{"me", "め", "メ", []string{"me"}, Mono},
	{"mo", "も", "モ", []string{"mo"}, Mono},
	{"ya", "や", "ヤ", []string{"ya"}, Mono},
	{"yu", "ゆ", "ユ", []string{"yu"}, Mono},
	{"yo", "よ", "ヨ", []string{"yo"}, Mono},
	{"ra", "ら", "ラ", []string{"ra"}, Mono},
	{"ri", "り", "リ", []string{"ri"}, Mono},
	{"ru", "る", "ル", []string{"ru"}, Mono},
	{"re", "れ", "レ", []string{"re"}, Mono},
	{"ro", "ろ", "ロ", []string{"ro"}, Mono},
	{"wa", "わ", "ワ", []string{"wa"}, Mono},
	{"wo", "を", "ヲ", []string{"wo"}, Mono},
	{"n", "ん", "ン", []string{"n"}, Mono},
	{"ga", "が", "ガ", []string{"ga"}, Dia},
	{"gi", "ぎ", "ギ", []string{"gi"}, Dia},
	{"gu", "ぐ", "グ", []string{"gu"}, Dia},
	{"ge", "げ", "ゲ", []string{"ge"}, Dia},
	{"go", "ご", "ゴ", []string{"go"}, Dia},
	{"za", "ざ", "ザ", []string{"za"}, Dia},
	{"ji", "じ", "ジ", []string{"ji", "zi"}, Dia},
	{"zu", "ず", "ズ", []string{"zu"}, Dia},
	{"ze", "ぜ", "ゼ", []string{"ze"}, Dia},
	{"zo", "ぞ", "ゾ", []string{"zo"}, Dia},
	{"da", "だ", "ダ", []string{"da"}, Dia},
	{"di", "ぢ", "ヂ", []string{"di", "dzi", "ji", "zi"}, Dia},
	{"du", "づ", "ヅ", []string{"du", "dzu", "zu"}, Dia},
	{"de", "で", "デ", []string{"de"}, Dia},
	{"do", "ど", "ド", []string{"do"}, Dia},
	{"ba", "ば", "バ", []string{"ba"}, Dia},
	{"bi", "び", "ビ", []string{"bi"}, Dia},
	{"bu", "ぶ", "ブ", []string{"bu"}, Dia},
	{"be", "べ", "ベ", []string{"be"}, Dia},
	{"bo", "ぼ", "ボ", []string{"bo"}, Dia},
	{"pa", "ぱ", "パ", []string{"pa"}, Dia},
	{"pi", "ぴ", "ピ", []string{"pi"}, Dia},
	{"pu", "ぷ", "プ", []string{"pu"}, Dia},
	{"pe", "ぺ", "ペ", []string{"pe"}, Dia},
	{"po", "ぽ", "ポ", []string{"po"}, Dia},
	{"kya", "きゃ", "キャ", []string{"kya"}, Di},
	{"kyu", "きゅ", "キュ", []string{"kyu"}, Di},
	{"kyo", "きょ", "キョ", []string{"kyo"}, Di},
	{"sha", "しゃ", "シャ", []string{"sha", "sya"}, Di},
	{"shu", "しゅ", "シュ", []string{"shu", "syu"}, Di},
	{"sho", "しょ", "ショ", []string{"sho", "syo"}, Di},
	{"cha", "ちゃ", "チャ", []string{"cha", "cya", "tya"}, Di},
	{"chu", "ちゅ", "チュ", []string{"chu", "cyu", "tyu"}, Di},
	{"cho", "ちょ", "チョ", []string{"cho", "cyo", "tyo"}, Di},
	{"nya", "にゃ", "ニャ", []string{"nya"}, Di},
	{"nyu", "にゅ", "ニュ", []string{"nyu"}, Di},
	{"nyo", "にょ", "ニョ", []string{"nyo"}, Di},
	{"hya", "ひゃ", "ヒャ", []string{"hya"}, Di},
	{"hyu", "ひゅ", "ヒュ", []string{"hyu"}, Di},
	{"hyo", "ひょ", "ヒョ", []string{"hyo"}, Di},
	{"mya", "みゃ", "ミャ", []string{"mya"}, Di},
	{"myu", "みゅ", "ミュ", []string{"myu"}, Di},
	{"myo", "みょ", "ミョ", []string{"myo"}, Di},
	{"rya", "りゃ", "リャ", []string{"rya"}, Di},
	{"ryu", "りゅ", "リュ", []string{"ryu"}, Di},
	{"ryo", "りょ", "リョ", []string{"ryo"}, Di},
	{"gya", "ぎゃ", "ギャ", []string{"gya"}, DiDia},
	{"gyu", "ぎゅ", "ギュ", []string{"gyu"}, DiDia},
	{"gyo", "ぎょ", "ギョ", []string{"gyo"}, DiDia},
	{"ja", "じゃ", "ジャ", []string{"ja", "jya", "zya"}, DiDia},
	{"ju", "じゅ", "ジュ", []string{"ju", "jyu", "zyu"}, DiDia},
	{"jo", "じょ", "ジョ", []string{"jo", "jyo", "zyo"}, DiDia},
	{"dya", "ぢゃ", "ヂャ", []string{"dya", "dzya", "ja", "jya", "zya"}, DiDia},
	{"dyu", "ぢゅ", "ヂュ", []string{"dyu", "dzyu", "ju", "jyu", "zyu"}, DiDia},
	{"dyo", "ぢょ", "ヂョ", []string{"dyo", "dzyo", "jo", "jyo", "zyo"}, DiDia},
	{"bya", "びゃ", "ビャ", []string{"bya"}, DiDia},
	{"byu", "びゅ", "ビュ", []string{"byu"}, DiDia},
	{"byo", "びょ", "ビョ", []string{"byo"}, DiDia},
	{"pya", "ぴゃ", "ピャ", []string{"pya"}, DiDia},
	{"pyu", "ぴゅ", "ピュ", []string{"pyu"}, DiDia},
	{"pyo", "ぴょ", "ピョ", []string{"pyo"}, DiDia},
}
