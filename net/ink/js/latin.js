"use strict";
/*
	Clive's Plan 9 like latin composition.
	kmap.islatin(str): prefix or full latin sequence
	kmap.latin(str): str for a full latin sequence.
*/

function KMap() {
	var m = new Map();
	var addseq = function (seq, s, nb) {
		var k = ""
		for(var i = 0; i < seq.length; i++) {
			k += seq[i];
			if (i == seq.length-1) {
				m.set(k, {s, nb});
			} else {
				m.set(k, true);
			}
		}
	};

	// return true if this is a full latin sequence or a prefix.
	this.islatin = function(s) {
		if (m.get(s)) {
			return true;
		}
		return false;
	};

	// return the rune (string) for a full latin sequence or undefined.
	this.latin = function(s) {
		var r = m.get(s);
		if (!r || r === true) {
			return undefined;
		}
		return r.s;
	};

	addseq("  ",	"␣", 0x2423);
	addseq(" i",	"ı", 0x0131);
	addseq("!~-",	"≄", 0x2244);
	addseq("!~=",	"≇", 0x2247);
	addseq("!~~",	"≉", 0x2249);
	addseq("!!",	"¡", 0x00a1);
	addseq("!<",	"≮", 0x226e);
	addseq("!=",	"≠", 0x2260);
	addseq("!>",	"≯", 0x226f);
	addseq("!?",	"‽", 0x203d);
	addseq("!b",	"⊄", 0x2284);
	addseq("!m",	"∉", 0x2209);
	addseq("!p",	"⊅", 0x2285);
	addseq("\"*I",	"Ϊ", 0x03aa);
	addseq("\"*U",	"Ϋ", 0x03ab);
	addseq("\"*i",	"ϊ", 0x03ca);
	addseq("\"*u",	"ϋ", 0x03cb);
	addseq("\"\"",	"¨", 0x00a8);
	addseq("\"A",	"Ä", 0x00c4);
	addseq("\"E",	"Ë", 0x00cb);
	addseq("\"I",	"Ï", 0x00cf);
	addseq("\"O",	"Ö", 0x00d6);
	addseq("\"U",	"Ü", 0x00dc);
	addseq("\"Y",	"Ÿ", 0x0178);
	addseq("\"a",	"ä", 0x00e4);
	addseq("\"e",	"ë", 0x00eb);
	addseq("\"i",	"ï", 0x00ef);
	addseq("\"o",	"ö", 0x00f6);
	addseq("\"u",	"ü", 0x00fc);
	addseq("\"y",	"ÿ", 0x00ff);
	addseq("$*f",	"ϕ", 0x03d5);
	addseq("$*h",	"ϑ", 0x03d1);
	addseq("$*k",	"ϰ", 0x03f0);
	addseq("$B",	"ℬ", 0x212c);
	addseq("$E",	"ℰ", 0x2130);
	addseq("$F",	"ℱ", 0x2131);
	addseq("$H",	"ℋ", 0x210b);
	addseq("$I",	"ℐ", 0x2110);
	addseq("$L",	"ℒ", 0x2112);
	addseq("$M",	"ℳ", 0x2133);
	addseq("$R",	"ℛ", 0x211b);
	addseq("$V",	"Ʋ", 0x01b2);
	addseq("$a",	"ɑ", 0x0251);
	addseq("$e",	"ℯ", 0x212f);
	addseq("$f",	"ƒ", 0x0192);
	addseq("$g",	"ℊ", 0x210a);
	addseq("$l",	"ℓ", 0x2113);
	addseq("$o",	"ℴ", 0x2134);
	addseq("$p",	"℘", 0x2118);
	addseq("$v",	"ʋ", 0x028b);
	addseq("'\"U",	"Ǘ", 0x01d7);
	addseq("'\"u",	"ǘ", 0x01d8);
	addseq("''",	"´", 0x00b4);
	addseq("'A",	"Á", 0x00c1);
	addseq("'C",	"Ć", 0x0106);
	addseq("'E",	"É", 0x00c9);
	addseq("'I",	"Í", 0x00cd);
	addseq("'L",	"Ĺ", 0x0139);
	addseq("'N",	"Ń", 0x0143);
	addseq("'O",	"Ó", 0x00d3);
	addseq("'R",	"Ŕ", 0x0154);
	addseq("'S",	"Ś", 0x015a);
	addseq("'U",	"Ú", 0x00da);
	addseq("'Y",	"Ý", 0x00dd);
	addseq("'Z",	"Ź", 0x0179);
	addseq("'a",	"á", 0x00e1);
	addseq("'c",	"ć", 0x0107);
	addseq("'e",	"é", 0x00e9);
	addseq("'g",	"ģ", 0x0123);
	addseq("'i",	"í", 0x00ed);
	addseq("'l",	"ĺ", 0x013a);
	addseq("'n",	"ń", 0x0144);
	addseq("'o",	"ó", 0x00f3);
	addseq("'r",	"ŕ", 0x0155);
	addseq("'s",	"ś", 0x015b);
	addseq("'u",	"ú", 0x00fa);
	addseq("'y",	"ý", 0x00fd);
	addseq("'z",	"ź", 0x017a);
	addseq("**",	"∗", 0x2217);
	addseq("*A",	"Α", 0x0391);
	addseq("*B",	"Β", 0x0392);
	addseq("*C",	"Ξ", 0x039e);
	addseq("*D",	"Δ", 0x0394);
	addseq("*E",	"Ε", 0x0395);
	addseq("*F",	"Φ", 0x03a6);
	addseq("*G",	"Γ", 0x0393);
	addseq("*H",	"Θ", 0x0398);
	addseq("*I",	"Ι", 0x0399);
	addseq("*K",	"Κ", 0x039a);
	addseq("*L",	"Λ", 0x039b);
	addseq("*M",	"Μ", 0x039c);
	addseq("*N",	"Ν", 0x039d);
	addseq("*O",	"Ο", 0x039f);
	addseq("*P",	"Π", 0x03a0);
	addseq("*Q",	"Ψ", 0x03a8);
	addseq("*R",	"Ρ", 0x03a1);
	addseq("*S",	"Σ", 0x03a3);
	addseq("*T",	"Τ", 0x03a4);
	addseq("*U",	"Υ", 0x03a5);
	addseq("*W",	"Ω", 0x03a9);
	addseq("*X",	"Χ", 0x03a7);
	addseq("*Y",	"Η", 0x0397);
	addseq("*Z",	"Ζ", 0x0396);
	addseq("*a",	"α", 0x03b1);
	addseq("*b",	"β", 0x03b2);
	addseq("*c",	"ξ", 0x03be);
	addseq("*d",	"δ", 0x03b4);
	addseq("*e",	"ε", 0x03b5);
	addseq("*f",	"φ", 0x03c6);
	addseq("*g",	"γ", 0x03b3);
	addseq("*h",	"θ", 0x03b8);
	addseq("*i",	"ι", 0x03b9);
	addseq("*k",	"κ", 0x03ba);
	addseq("*l",	"λ", 0x03bb);
	addseq("*m",	"μ", 0x03bc);
	addseq("*n",	"ν", 0x03bd);
	addseq("*o",	"ο", 0x03bf);
	addseq("*p",	"π", 0x03c0);
	addseq("*q",	"ψ", 0x03c8);
	addseq("*r",	"ρ", 0x03c1);
	addseq("*s",	"σ", 0x03c3);
	addseq("*t",	"τ", 0x03c4);
	addseq("*u",	"υ", 0x03c5);
	addseq("*w",	"ω", 0x03c9);
	addseq("*x",	"χ", 0x03c7);
	addseq("*y",	"η", 0x03b7);
	addseq("*z",	"ζ", 0x03b6);
	addseq("+-",	"±", 0x00b1);
	addseq("+O",	"⊕", 0x2295);
	addseq(",,",	"¸", 0x00b8);
	addseq(",A",	"Ą", 0x0104);
	addseq(",C",	"Ç", 0x00c7);
	addseq(",E",	"Ę", 0x0118);
	addseq(",G",	"Ģ", 0x0122);
	addseq(",I",	"Į", 0x012e);
	addseq(",K",	"Ķ", 0x0136);
	addseq(",L",	"Ļ", 0x013b);
	addseq(",N",	"Ņ", 0x0145);
	addseq(",O",	"Ǫ", 0x01ea);
	addseq(",R",	"Ŗ", 0x0156);
	addseq(",S",	"Ş", 0x015e);
	addseq(",T",	"Ţ", 0x0162);
	addseq(",U",	"Ų", 0x0172);
	addseq(",a",	"ą", 0x0105);
	addseq(",c",	"ç", 0x00e7);
	addseq(",e",	"ę", 0x0119);
	addseq(",g",	"ģ", 0x0123);
	addseq(",i",	"į", 0x012f);
	addseq(",k",	"ķ", 0x0137);
	addseq(",l",	"ļ", 0x013c);
	addseq(",n",	"ņ", 0x0146);
	addseq(",o",	"ǫ", 0x01eb);
	addseq(",r",	"ŗ", 0x0157);
	addseq(",s",	"ş", 0x015f);
	addseq(",t",	"ţ", 0x0163);
	addseq(",u",	"ų", 0x0173);
	addseq("-*l",	"ƛ", 0x019b);
	addseq("-+",	"∓", 0x2213);
	addseq("--",	"­", 0x00ad);
	addseq("-2",	"ƻ", 0x01bb);
	addseq("-:",	"÷", 0x00f7);
	addseq("->",	"→", 0x2192);
	addseq("-D",	"Ð", 0x00d0);
	addseq("-G",	"Ǥ", 0x01e4);
	addseq("-H",	"Ħ", 0x0126);
	addseq("-I",	"Ɨ", 0x0197);
	addseq("-L",	"Ł", 0x0141);
	addseq("-O",	"⊖", 0x2296);
	addseq("-T",	"Ŧ", 0x0166);
	addseq("-Z",	"Ƶ", 0x01b5);
	addseq("-b",	"ƀ", 0x0180);
	addseq("-d",	"ð", 0x00f0);
	addseq("-g",	"ǥ", 0x01e5);
	addseq("-h",	"ℏ", 0x210f);
	addseq("-i",	"ɨ", 0x0268);
	addseq("-l",	"ł", 0x0142);
	addseq("-t",	"ŧ", 0x0167);
	addseq("-u",	"ʉ", 0x0289);
	addseq("-z",	"ƶ", 0x01b6);
	addseq("-~",	"≂", 0x2242);
	addseq("..",	"·", 0x00b7);
	addseq(".C",	"Ċ", 0x010a);
	addseq(".E",	"Ė", 0x0116);
	addseq(".G",	"Ġ", 0x0120);
	addseq(".I",	"İ", 0x0130);
	addseq(".L",	"Ŀ", 0x013f);
	addseq(".O",	"⊙", 0x2299);
	addseq(".Z",	"Ż", 0x017b);
	addseq(".c",	"ċ", 0x010b);
	addseq(".e",	"ė", 0x0117);
	addseq(".g",	"ġ", 0x0121);
	addseq(".l",	"ŀ", 0x0140);
	addseq(".z",	"ż", 0x017c);
	addseq("/O",	"Ø", 0x00d8);
	addseq("/o",	"ø", 0x00f8);
	addseq("1.",	"․", 0x2024);
	addseq("12",	"½", 0x00bd);
	addseq("13",	"⅓", 0x2153);
	addseq("14",	"¼", 0x00bc);
	addseq("15",	"⅕", 0x2155);
	addseq("16",	"⅙", 0x2159);
	addseq("18",	"⅛", 0x215b);
	addseq("2-",	"ƻ", 0x01bb);
	addseq("2.",	"‥", 0x2025);
	addseq("23",	"⅔", 0x2154);
	addseq("25",	"⅖", 0x2156);
	addseq("3.",	"…", 0x2026);
	addseq("34",	"¾", 0x00be);
	addseq("35",	"⅗", 0x2157);
	addseq("38",	"⅜", 0x215c);
	addseq("45",	"⅘", 0x2158);
	addseq("56",	"⅚", 0x215a);
	addseq("58",	"⅝", 0x215d);
	addseq("78",	"⅞", 0x215e);
	addseq(":(",	"☹", 0x2639);
	addseq(":)",	"☺", 0x263a);
	addseq(":-",	"÷", 0x00f7);
	addseq(":=",	"≔", 0x2254);
	addseq("<!=",	"≨", 0x2268);
	addseq("<!~",	"⋦", 0x22e6);
	addseq("<-",	"←", 0x2190);
	addseq("<<",	"«", 0x00ab);
	addseq("<=",	"≤", 0x2264);
	addseq("<>",	"≶", 0x2276);
	addseq("<~",	"≲", 0x2272);
	addseq("=:",	"≕", 0x2255);
	addseq("=<",	"⋜", 0x22dc);
	addseq("==",	"≡", 0x2261);
	addseq("=>",	"⋝", 0x22dd);
	addseq("=O",	"⊜", 0x229c);
	addseq("=V",	"⇒", 0x21d2);
	addseq(">!=",	"≩", 0x2269);
	addseq(">!~",	"⋧", 0x22e7);
	addseq("><",	"≷", 0x2277);
	addseq(">=",	"≥", 0x2265);
	addseq(">>",	"»", 0x00bb);
	addseq(">~",	"≳", 0x2273);
	addseq("?!",	"‽", 0x203d);
	addseq("??",	"¿", 0x00bf);
	addseq("@''",	"ъ", 0x044a);
	addseq("@@'",	"ь", 0x044c);
	addseq("@@E",	"Е", 0x0415);
	addseq("@@K",	"К", 0x041a);
	addseq("@@S",	"С", 0x0421);
	addseq("@@T",	"Т", 0x0422);
	addseq("@@Y",	"Ы", 0x042b);
	addseq("@@Z",	"З", 0x0417);
	addseq("@@e",	"е", 0x0435);
	addseq("@@k",	"к", 0x043a);
	addseq("@@s",	"с", 0x0441);
	addseq("@@t",	"т", 0x0442);
	addseq("@@y",	"ы", 0x044b);
	addseq("@@z",	"з", 0x0437);
	addseq("@CH",	"Ч", 0x0427);
	addseq("@Ch",	"Ч", 0x0427);
	addseq("@EH",	"Э", 0x042d);
	addseq("@Eh",	"Э", 0x042d);
	addseq("@KH",	"Х", 0x0425);
	addseq("@Kh",	"Х", 0x0425);
	addseq("@SC",	"Щ", 0x0429);
	addseq("@SH",	"Ш", 0x0428);
	addseq("@Sc",	"Щ", 0x0429);
	addseq("@Sh",	"Ш", 0x0428);
	addseq("@TS",	"Ц", 0x0426);
	addseq("@Ts",	"Ц", 0x0426);
	addseq("@YA",	"Я", 0x042f);
	addseq("@YE",	"Е", 0x0415);
	addseq("@YO",	"Ё", 0x0401);
	addseq("@YU",	"Ю", 0x042e);
	addseq("@Ya",	"Я", 0x042f);
	addseq("@Ye",	"Е", 0x0415);
	addseq("@Yo",	"Ё", 0x0401);
	addseq("@Yu",	"Ю", 0x042e);
	addseq("@ZH",	"Ж", 0x0416);
	addseq("@Zh",	"Ж", 0x0416);
	addseq("@ch",	"ч", 0x0447);
	addseq("@eh",	"э", 0x044d);
	addseq("@kh",	"х", 0x0445);
	addseq("@sc",	"щ", 0x0449);
	addseq("@sh",	"ш", 0x0448);
	addseq("@ts",	"ц", 0x0446);
	addseq("@ya",	"я", 0x044f);
	addseq("@ye",	"е", 0x0435);
	addseq("@yo",	"ё", 0x0451);
	addseq("@yu",	"ю", 0x044e);
	addseq("@zh",	"ж", 0x0436);
	addseq("@A",	"А", 0x0410);
	addseq("@B",	"Б", 0x0411);
	addseq("@D",	"Д", 0x0414);
	addseq("@F",	"Ф", 0x0424);
	addseq("@G",	"Г", 0x0413);
	addseq("@I",	"И", 0x0418);
	addseq("@J",	"Й", 0x0419);
	addseq("@L",	"Л", 0x041b);
	addseq("@M",	"М", 0x041c);
	addseq("@N",	"Н", 0x041d);
	addseq("@O",	"О", 0x041e);
	addseq("@P",	"П", 0x041f);
	addseq("@R",	"Р", 0x0420);
	addseq("@U",	"У", 0x0423);
	addseq("@V",	"В", 0x0412);
	addseq("@X",	"Х", 0x0425);
	addseq("@a",	"а", 0x0430);
	addseq("@b",	"б", 0x0431);
	addseq("@d",	"д", 0x0434);
	addseq("@f",	"ф", 0x0444);
	addseq("@g",	"г", 0x0433);
	addseq("@i",	"и", 0x0438);
	addseq("@j",	"й", 0x0439);
	addseq("@l",	"л", 0x043b);
	addseq("@m",	"м", 0x043c);
	addseq("@n",	"н", 0x043d);
	addseq("@o",	"о", 0x043e);
	addseq("@p",	"п", 0x043f);
	addseq("@r",	"р", 0x0440);
	addseq("@u",	"у", 0x0443);
	addseq("@v",	"в", 0x0432);
	addseq("@x",	"х", 0x0445);
	addseq("AE",	"Æ", 0x00c6);
	addseq("CA",	"⋂", 0x22c2);
	addseq("CC",	"ℂ", 0x2102);
	addseq("CU",	"⋃", 0x22c3);
	addseq("DvZ",	"Ǆ", 0x01c4);
	addseq("Dvz",	"ǅ", 0x01c5);
	addseq("D-",	"Ð", 0x00d0);
	addseq("De",	"∆", 0x2206);
	addseq("G-",	"Ǥ", 0x01e4);
	addseq("H-",	"Ħ", 0x0126);
	addseq("HH",	"ℍ", 0x210d);
	addseq("I-",	"Ɨ", 0x0197);
	addseq("IJ",	"Ĳ", 0x0132);
	addseq("L&",	"⋀", 0x22c0);
	addseq("L-",	"Ł", 0x0141);
	addseq("LJ",	"Ǉ", 0x01c7);
	addseq("Lj",	"ǈ", 0x01c8);
	addseq("L|",	"⋁", 0x22c1);
	addseq("M#",	"♮", 0x266e);
	addseq("M4",	"♩", 0x2669);
	addseq("M8",	"♪", 0x266a);
	addseq("Mb",	"♭", 0x266d);
	addseq("Ms",	"♯", 0x266f);
	addseq("NJ",	"Ǌ", 0x01ca);
	addseq("NN",	"ℕ", 0x2115);
	addseq("Nj",	"ǋ", 0x01cb);
	addseq("O*",	"⊛", 0x229b);
	addseq("O+",	"⊕", 0x2295);
	addseq("O-",	"⊖", 0x2296);
	addseq("O.",	"⊙", 0x2299);
	addseq("O/",	"⊘", 0x2298);
	addseq("O=",	"⊜", 0x229c);
	addseq("OE",	"Œ", 0x0152);
	addseq("OI",	"Ƣ", 0x01a2);
	addseq("Oc",	"©", 0x00a9);
	addseq("Oo",	"⊚", 0x229a);
	addseq("Op",	"℗", 0x2117);
	addseq("Or",	"®", 0x00ae);
	addseq("Ox",	"⊗", 0x2297);
	addseq("PP",	"ℙ", 0x2119);
	addseq("QQ",	"ℚ", 0x211a);
	addseq("RR",	"ℝ", 0x211d);
	addseq("S1",	"¹", 0x00b9);
	addseq("S2",	"²", 0x00b2);
	addseq("S3",	"³", 0x00b3);
	addseq("SS",	"§", 0x00a7);
	addseq("T-",	"Ŧ", 0x0166);
	addseq("Tu",	"⊨", 0x22a8);
	addseq("V=",	"⇐", 0x21d0);
	addseq("YR",	"Ʀ", 0x01a6);
	addseq("Z-",	"Ƶ", 0x01b5);
	addseq("ZA",	"", 0xf015);
	addseq("ZC",	"", 0xf017);
	addseq("ZS",	"", 0xf016);
	addseq("ZZ",	"ℤ", 0x2124);
	addseq("^A",	"Â", 0x00c2);
	addseq("^C",	"Ĉ", 0x0108);
	addseq("^E",	"Ê", 0x00ca);
	addseq("^G",	"Ĝ", 0x011c);
	addseq("^H",	"Ĥ", 0x0124);
	addseq("^I",	"Î", 0x00ce);
	addseq("^J",	"Ĵ", 0x0134);
	addseq("^O",	"Ô", 0x00d4);
	addseq("^S",	"Ŝ", 0x015c);
	addseq("^U",	"Û", 0x00db);
	addseq("^W",	"Ŵ", 0x0174);
	addseq("^Y",	"Ŷ", 0x0176);
	addseq("^a",	"â", 0x00e2);
	addseq("^c",	"ĉ", 0x0109);
	addseq("^e",	"ê", 0x00ea);
	addseq("^g",	"ĝ", 0x011d);
	addseq("^h",	"ĥ", 0x0125);
	addseq("^i",	"î", 0x00ee);
	addseq("^j",	"ĵ", 0x0135);
	addseq("^o",	"ô", 0x00f4);
	addseq("^s",	"ŝ", 0x015d);
	addseq("^u",	"û", 0x00fb);
	addseq("^w",	"ŵ", 0x0175);
	addseq("^y",	"ŷ", 0x0177);
	addseq("_\"A",	"Ǟ", 0x01de);
	addseq("_\"U",	"Ǖ", 0x01d5);
	addseq("_\"a",	"ǟ", 0x01df);
	addseq("_\"u",	"ǖ", 0x01d6);
	addseq("_,O",	"Ǭ", 0x01ec);
	addseq("_,o",	"ǭ", 0x01ed);
	addseq("_.A",	"Ǡ", 0x01e0);
	addseq("_.a",	"ǡ", 0x01e1);
	addseq("_A",	"Ā", 0x0100);
	addseq("_E",	"Ē", 0x0112);
	addseq("_I",	"Ī", 0x012a);
	addseq("_O",	"Ō", 0x014c);
	addseq("_U",	"Ū", 0x016a);
	addseq("__",	"¯", 0x00af);
	addseq("_a",	"ā", 0x0101);
	addseq("_e",	"ē", 0x0113);
	addseq("_i",	"ī", 0x012b);
	addseq("_o",	"ō", 0x014d);
	addseq("_u",	"ū", 0x016b);
	addseq("`\"U",	"Ǜ", 0x01db);
	addseq("`\"u",	"ǜ", 0x01dc);
	addseq("`A",	"À", 0x00c0);
	addseq("`E",	"È", 0x00c8);
	addseq("`I",	"Ì", 0x00cc);
	addseq("`O",	"Ò", 0x00d2);
	addseq("`U",	"Ù", 0x00d9);
	addseq("`a",	"à", 0x00e0);
	addseq("`e",	"è", 0x00e8);
	addseq("`i",	"ì", 0x00ec);
	addseq("`o",	"ò", 0x00f2);
	addseq("`u",	"ù", 0x00f9);
	addseq("ab",	"↔", 0x2194);
	addseq("ae",	"æ", 0x00e6);
	addseq("an",	"∠", 0x2220);
	addseq("b(",	"₍", 0x208d);
	addseq("b)",	"₎", 0x208e);
	addseq("b+",	"₊", 0x208a);
	addseq("b-",	"₋", 0x208b);
	addseq("b0",	"₀", 0x2080);
	addseq("b1",	"₁", 0x2081);
	addseq("b2",	"₂", 0x2082);
	addseq("b3",	"₃", 0x2083);
	addseq("b4",	"₄", 0x2084);
	addseq("b5",	"₅", 0x2085);
	addseq("b6",	"₆", 0x2086);
	addseq("b7",	"₇", 0x2087);
	addseq("b8",	"₈", 0x2088);
	addseq("b9",	"₉", 0x2089);
	addseq("b=",	"₌", 0x208c);
	addseq("bb",	"♝", 0x265d);
	addseq("bk",	"♚", 0x265a);
	addseq("bn",	"♞", 0x265e);
	addseq("bp",	"♟", 0x265f);
	addseq("bq",	"♛", 0x265b);
	addseq("br",	"♜", 0x265c);
	addseq("bu",	"•", 0x2022);
	addseq("c$",	"¢", 0x00a2);
	addseq("cO",	"©", 0x00a9);
	addseq("ca",	"∩", 0x2229);
	addseq("cg",	"≅", 0x2245);
	addseq("cu",	"∪", 0x222a);
	addseq("dvz",	"ǆ", 0x01c6);
	addseq("d-",	"ð", 0x00f0);
	addseq("da",	"↓", 0x2193);
	addseq("dd",	"‡", 0x2021);
	addseq("de",	"°", 0x00b0);
	addseq("dg",	"†", 0x2020);
	addseq("dz",	"ʣ", 0x02a3);
	addseq("e$",	"€", 0x20ac);
	addseq("el",	"⋯", 0x22ef);
	addseq("em",	"—", 0x2014);
	addseq("en",	"–", 0x2013);
	addseq("es",	"∅", 0x2205);
	addseq("fa",	"∀", 0x2200);
	addseq("g$",	"¤", 0x00a4);
	addseq("g-",	"ǥ", 0x01e5);
	addseq("gr",	"∇", 0x2207);
	addseq("h-",	"ℏ", 0x210f);
	addseq("hv",	"ƕ", 0x0195);
	addseq("i-",	"ɨ", 0x0268);
	addseq("ib",	"⊆", 0x2286);
	addseq("if",	"∞", 0x221e);
	addseq("ij",	"ĳ", 0x0133);
	addseq("ip",	"⊇", 0x2287);
	addseq("is",	"∫", 0x222b);
	addseq("l\"",	"“", 0x201c);
	addseq("l$",	"£", 0x00a3);
	addseq("l&",	"∧", 0x2227);
	addseq("l'",	"‘", 0x2018);
	addseq("l-",	"ł", 0x0142);
	addseq("lj",	"ǉ", 0x01c9);
	addseq("lz",	"⋄", 0x22c4);
	addseq("l|",	"∨", 0x2228);
	addseq("mi",	"µ", 0x00b5);
	addseq("mo",	"∈", 0x2208);
	addseq("mu",	"×", 0x00d7);
	addseq("nj",	"ǌ", 0x01cc);
	addseq("no",	"¬", 0x00ac);
	addseq("oA",	"Å", 0x00c5);
	addseq("oO",	"⊚", 0x229a);
	addseq("oU",	"Ů", 0x016e);
	addseq("oa",	"å", 0x00e5);
	addseq("oe",	"œ", 0x0153);
	addseq("oi",	"ƣ", 0x01a3);
	addseq("ou",	"ů", 0x016f);
	addseq("pO",	"℗", 0x2117);
	addseq("pd",	"∂", 0x2202);
	addseq("pg",	"¶", 0x00b6);
	addseq("pr",	"∏", 0x220f);
	addseq("pt",	"∝", 0x221d);
	addseq("r\"",	"”", 0x201d);
	addseq("r'",	"’", 0x2019);
	addseq("rO",	"®", 0x00ae);
	addseq("s(",	"⁽", 0x207d);
	addseq("s)",	"⁾", 0x207e);
	addseq("s+",	"⁺", 0x207a);
	addseq("s-",	"⁻", 0x207b);
	addseq("s0",	"⁰", 0x2070);
	addseq("s1",	"¹", 0x00b9);
	addseq("s2",	"²", 0x00b2);
	addseq("s3",	"³", 0x00b3);
	addseq("s4",	"⁴", 0x2074);
	addseq("s5",	"⁵", 0x2075);
	addseq("s6",	"⁶", 0x2076);
	addseq("s7",	"⁷", 0x2077);
	addseq("s8",	"⁸", 0x2078);
	addseq("s9",	"⁹", 0x2079);
	addseq("s=",	"⁼", 0x207c);
	addseq("sa",	"ª", 0x00aa);
	addseq("sb",	"⊂", 0x2282);
	addseq("si",	"ⁱ", 0x2071);
	addseq("sn",	"ⁿ", 0x207f);
	addseq("so",	"º", 0x00ba);
	addseq("sp",	"⊃", 0x2283);
	addseq("sr",	"√", 0x221a);
	addseq("ss",	"ß", 0x00df);
	addseq("st",	"∍", 0x220d);
	addseq("su",	"∑", 0x2211);
	addseq("t-",	"ŧ", 0x0167);
	addseq("te",	"∃", 0x2203);
	addseq("tf",	"∴", 0x2234);
	addseq("tm",	"™", 0x2122);
	addseq("ts",	"ς", 0x03c2);
	addseq("tu",	"⊢", 0x22a2);
	addseq("u-",	"ʉ", 0x0289);
	addseq("uA",	"Ă", 0x0102);
	addseq("uE",	"Ĕ", 0x0114);
	addseq("uG",	"Ğ", 0x011e);
	addseq("uI",	"Ĭ", 0x012c);
	addseq("uO",	"Ŏ", 0x014e);
	addseq("uU",	"Ŭ", 0x016c);
	addseq("ua",	"↑", 0x2191);
	addseq("ue",	"ĕ", 0x0115);
	addseq("ug",	"ğ", 0x011f);
	addseq("ui",	"ĭ", 0x012d);
	addseq("uo",	"ŏ", 0x014f);
	addseq("uu",	"ŭ", 0x016d);
	addseq("v\"U",	"Ǚ", 0x01d9);
	addseq("v\"u",	"ǚ", 0x01da);
	addseq("vA",	"Ǎ", 0x01cd);
	addseq("vC",	"Č", 0x010c);
	addseq("vD",	"Ď", 0x010e);
	addseq("vE",	"Ě", 0x011a);
	addseq("vG",	"Ǧ", 0x01e6);
	addseq("vI",	"Ǐ", 0x01cf);
	addseq("vK",	"Ǩ", 0x01e8);
	addseq("vL",	"Ľ", 0x013d);
	addseq("vN",	"Ň", 0x0147);
	addseq("vO",	"Ǒ", 0x01d1);
	addseq("vR",	"Ř", 0x0158);
	addseq("vS",	"Š", 0x0160);
	addseq("vT",	"Ť", 0x0164);
	addseq("vU",	"Ǔ", 0x01d3);
	addseq("vZ",	"Ž", 0x017d);
	addseq("va",	"ǎ", 0x01ce);
	addseq("vc",	"č", 0x010d);
	addseq("vd",	"ď", 0x010f);
	addseq("ve",	"ě", 0x011b);
	addseq("vg",	"ǧ", 0x01e7);
	addseq("vi",	"ǐ", 0x01d0);
	addseq("vj",	"ǰ", 0x01f0);
	addseq("vk",	"ǩ", 0x01e9);
	addseq("vl",	"ľ", 0x013e);
	addseq("vn",	"ň", 0x0148);
	addseq("vo",	"ǒ", 0x01d2);
	addseq("vr",	"ř", 0x0159);
	addseq("vs",	"š", 0x0161);
	addseq("vt",	"ť", 0x0165);
	addseq("vu",	"ǔ", 0x01d4);
	addseq("vz",	"ž", 0x017e);
	addseq("wb",	"♗", 0x2657);
	addseq("wk",	"♔", 0x2654);
	addseq("wn",	"♘", 0x2658);
	addseq("wp",	"♙", 0x2659);
	addseq("wq",	"♕", 0x2655);
	addseq("wr",	"♖", 0x2656);
	addseq("xO",	"⊗", 0x2297);
	addseq("y$",	"¥", 0x00a5);
	addseq("z-",	"ƶ", 0x01b6);
	addseq("|P",	"Þ", 0x00de);
	addseq("|p",	"þ", 0x00fe);
	addseq("||",	"¦", 0x00a6);
	addseq("~!=",	"≆", 0x2246);
	addseq("~-",	"≃", 0x2243);
	addseq("~=",	"≅", 0x2245);
	addseq("~A",	"Ã", 0x00c3);
	addseq("~I",	"Ĩ", 0x0128);
	addseq("~N",	"Ñ", 0x00d1);
	addseq("~O",	"Õ", 0x00d5);
	addseq("~U",	"Ũ", 0x0168);
	addseq("~a",	"ã", 0x00e3);
	addseq("~i",	"ĩ", 0x0129);
	addseq("~n",	"ñ", 0x00f1);
	addseq("~o",	"õ", 0x00f5);
	addseq("~u",	"ũ", 0x0169);
	addseq("~~",	"≈", 0x2248);
}

var kmap = new KMap();