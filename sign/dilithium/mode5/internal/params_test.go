package internal

import (
	"testing"

	"github.com/JI-0/circl-ed448/sign/dilithium/internal/common"
)

// Tests specific to the current mode

func TestVectorDeriveUniformLeqEta(t *testing.T) {
	var p common.Poly
	var seed [64]byte
	p2 := common.Poly{
		8380415, 1, 8380415, 8380416, 1, 8380415, 8380416, 8380415,
		8380415, 8380416, 8380416, 2, 2, 1, 8380415, 8380415,
		8380416, 1, 8380416, 8380415, 2, 0, 0, 1, 1, 2, 2, 8380415,
		0, 8380416, 8380416, 8380416, 8380415, 1, 2, 0, 1, 8380415,
		0, 1, 8380415, 1, 0, 8380415, 2, 1, 2, 0, 1, 0, 8380416,
		1, 8380416, 1, 0, 1, 1, 0, 1, 8380416, 0, 0, 8380416,
		8380415, 8380416, 2, 0, 0, 8380415, 1, 1, 0, 0, 1, 8380415,
		1, 8380416, 1, 8380415, 8380416, 8380416, 8380415, 0, 1,
		8380415, 8380415, 1, 8380415, 0, 2, 2, 8380415, 1, 2,
		8380415, 8380415, 0, 2, 2, 1, 8380415, 8380416, 0, 8380415,
		2, 1, 8380415, 2, 2, 8380416, 8380416, 0, 8380416, 0, 2,
		8380416, 1, 8380415, 8380416, 8380415, 1, 8380416, 8380416,
		2, 2, 0, 0, 0, 8380415, 8380415, 2, 8380416, 2, 2, 8380415,
		8380415, 2, 2, 2, 8380415, 1, 2, 1, 2, 8380415, 0, 2, 1,
		8380415, 2, 8380415, 8380415, 8380416, 0, 8380416, 8380415,
		8380415, 8380416, 8380416, 2, 8380416, 2, 0, 0, 1, 1, 1,
		8380416, 0, 8380416, 8380416, 1, 1, 1, 0, 8380416, 2, 0,
		8380415, 8380415, 0, 0, 2, 8380416, 1, 0, 0, 8380415,
		8380415, 1, 0, 8380416, 1, 2, 8380415, 0, 8380416, 8380415,
		1, 1, 0, 1, 8380416, 8380415, 1, 0, 0, 8380416, 1, 0, 2,
		8380416, 2, 2, 0, 0, 1, 1, 2, 8380415, 2, 8380416, 8380416,
		2, 1, 2, 8380416, 8380415, 8380415, 8380415, 0, 8380416,
		1, 0, 2, 8380416, 2, 8380415, 8380415, 2, 2, 8380415,
		8380416, 0, 8380415, 8380415, 0, 2, 8380415, 1, 8380415,
		8380415, 1, 1, 8380416, 8380416,
	}
	for i := 0; i < 64; i++ {
		seed[i] = byte(i)
	}
	PolyDeriveUniformLeqEta(&p, &seed, 30000)
	p.Normalize()
	if p != p2 {
		t.Fatalf("%v != %v", p, p2)
	}
}

func TestVectorDeriveUniformLeGamma1(t *testing.T) {
	var p, p2 common.Poly
	var seed [64]byte
	p2 = common.Poly{
		8011853, 7949494, 172552, 263871, 8095275, 155369, 311506,
		8076900, 8307558, 8139232, 8041607, 448815, 380634, 180526,
		8165391, 101857, 8286792, 427645, 8098920, 7860396, 352757,
		8179719, 7954627, 7898860, 28800, 8129086, 111121, 8115657,
		8211418, 7943538, 259410, 7965184, 8232538, 7864584,
		7991749, 23725, 393449, 8344363, 8041712, 196742, 8187277,
		230211, 115522, 205750, 8332267, 8020968, 511882, 66518,
		8377952, 283731, 276156, 488847, 218386, 24973, 7960226,
		8019608, 8163770, 8099393, 8251752, 8055784, 438808, 408276,
		245718, 90648, 8179442, 377149, 66371, 8067974, 8165213,
		496174, 7959821, 8174846, 416247, 8334586, 8277522, 137692,
		8260481, 45327, 8078022, 8223800, 8070188, 8291718, 156021,
		516504, 8144827, 361012, 323861, 8315499, 8004848, 7906709,
		7913063, 230858, 311998, 8280928, 8347571, 8236825, 120069,
		412722, 476656, 372912, 8036734, 465145, 8275725, 8153834,
		411759, 412681, 72836, 8378216, 8305773, 8162477, 8293183,
		289061, 7900478, 8133091, 100678, 267462, 254283, 242941,
		8009771, 364316, 217523, 8026537, 7899325, 7863708, 211663,
		339314, 8133229, 8035753, 135557, 8245724, 7988629, 8042510,
		8012465, 386933, 8351229, 88508, 274815, 8293482, 216047,
		8232256, 8337777, 8305592, 7938394, 378619, 7942432,
		7961498, 360341, 265269, 8346169, 514971, 8255059, 406815,
		222421, 8344231, 464482, 94984, 8147964, 8242727, 8211462,
		7945005, 8167987, 8290153, 8355124, 303031, 180689, 97653,
		8032319, 263210, 684, 437628, 7983244, 359393, 8054335,
		223796, 8014878, 8066876, 335829, 467349, 105150, 326057,
		229928, 7934510, 26854, 8093051, 8162834, 8013975, 8122355,
		44783, 7969925, 465863, 8299023, 8155688, 8256445, 7975782,
		7892171, 8075999, 412728, 7858411, 480155, 7922893, 254722,
		381253, 8307390, 8040031, 280413, 8089206, 7869244, 8050145,
		8028110, 8020538, 8158686, 7875907, 7960483, 7998991,
		8317674, 52939, 416219, 501681, 231283, 8151233, 8241847,
		8224119, 454076, 8171231, 411693, 8324986, 447356, 400055,
		490491, 477035, 8055459, 158775, 383762, 8167063, 8076788,
		7956883, 309585, 111368, 8312360, 7992502, 8259793, 461240,
		7937002, 8198300, 7862862, 302423, 437299, 420919, 8359979,
		8191730, 7895992, 75500, 307359, 435102, 7873624, 457428,
	}
	for i := 0; i < 64; i++ {
		seed[i] = byte(i)
	}
	PolyDeriveUniformLeGamma1(&p, &seed, 30000)
	p.Normalize()
	if p != p2 {
		t.Fatalf("%v != %v", p, p2)
	}
}
