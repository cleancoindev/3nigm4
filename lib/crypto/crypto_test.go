//
// 3nigm4 crypto package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 06/03/2016
//

package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"testing"
)

var (
	plaintex = `This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted. 
	This is a test message, it will be encrypted.`
	key1 = "key12345a293bf93key12345a293bf93"
	key2 = "this.is.a.test.key"
	key3 = "R39eie93oe0903i9e£eoo093"
	salt = "i93ie93e"
)

// https://raw.githubusercontent.com/golang/crypto/master/openpgp/keys_test.go
const publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

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
lw1nIDx1uwkeAfXXcViFBg==
=Ifs7
-----END PGP PUBLIC KEY BLOCK-----`

// pub   1024R/7F98BBCE 2014-01-04
// uid   Golang Test (Private key password is 'golang') <golangtest@test.com>
// sub   1024R/5F34A320 2014-01-04
const privateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1
 
lQH+BFLHbYYBBADCjgKHmPmwBxI3c3DPVoSdu0+EJl/EsS2HEaN63dnLkGsMAs+4
32wsywmMrzKqCL40sbhJVYBcfe0chL+cry4O54DX7+gA0ZSVzFUN2EGocnkaHzyS
fuUtBdCTmoWZZAGFiBwlIS7aE/86SOyHksFo8LRC9W/GIWQS2PbcadvUywARAQAB
/gMDApJxOwcsfChBYCCmhOAvotKdYcy7nuG7dyGDBlpclLJtH/PaakKSE33NtEj4
1fyixQOdwApxvuQ2P0VX3pie/De1KpbeqXfnPLsmsXQwrRPOo38T5zeJ5ToWUGDC
Oia69ep3kmHbAW41EBH/uk/nMM91QUdl4mkYsc3dhVOXbmf0xyRoP/Afqha4UhdZ
0XKlIZP1a5+3NF/Q6dAVG0+FlO5Hcai8n98jW0id8Yf6zI+1gFGvYYKhlifkdJeK
Nf4YEvOXALEvaQqkcJOxEca+BmqsgCIFctJe9Bahx97Ep5hP7AH0aBmtZfmGmZwB
GYoevUtKa4ASVmK8RaddBvIjcrWsoAsYMpDGYaE0fcdtxsBf3uT1Q8IMsT+ZRjjV
TfvJ8aW14ZrLI98KdtXaOPZs91mML+3iw1c/1O/IEJfwxrUni2p/fDmCYU9eHR3u
Q0PwVR0MCUHI1fGuUoetW2gYIxfklvBtEFWW1BD6fCpCtERHb2xhbmcgVGVzdCAo
UHJpdmF0ZSBrZXkgcGFzc3dvcmQgaXMgJ2dvbGFuZycpIDxnb2xhbmd0ZXN0QHRl
c3QuY29tPoi4BBMBAgAiBQJSx22GAhsDBgsJCAcDAgYVCAIJCgsEFgIDAQIeAQIX
gAAKCRBVSiCHf5i7zqKJA/sFUM2TfL2VZKWC7E1N1wwZctB9Bf77SeAPSVpGCZ0c
iUYIFdwwGowKtjoDrsbYgPp+UGOyYMD6tGzWKaJrQQoDyaQqVVRhbNXB7Jz7JT2a
qKHD1t7cx5FfUzDMBNou3TOWHomDXyQGDAULAZnjaOj8/pDe6poxyBluSjMJUzfD
pp0B/gRSx22GAQQArUMDqkGng9Cppk73UBWBd7jhhbtk0eaRQh/goUHhKJerZ4LM
Q21IKyIX+GQbscDpccpXMI6eThXxrL+D8G4cNb4ewvT0zc20+T91ztgT9A/4Vifc
EPQCErTqY/oZphAzZM1p6sRenc22e42iT0Iibd5gCs2wnSNeUzybDcuQi2EAEQEA
Af4DAwKScTsHLHwoQWCYayWqio8purPTonYogZSN3QwaheS2Y0NE7skdLOvP97vi
Rh7BktS6Dkgu0T3D39+q0O6ZO7XErvTVoas1F0HXzId4tiIicmx4tYNyWI4NrSO7
6TQPz/bQe8ZN+plG5cgZowts6g6RSfQxoW21LrP8Lh+OEdcYwWf7BTukAYmD3oq9
RxdfYI7hnbVGFdOqQUQNcxZkbdrsF9ITjQb/KRln5/99E1Kp1D45VpPOs7NT3orA
mnfSslJXVNm1uK6FDBX2iUe3JaAmgh+RLGXQXRZKJW4DGDTyYdwR4hO8cYix2+8z
+XuwdVDPKBnzKn190m6xpdLyvKfj1BQhX14NShPQZ3QJiMU0k4Js23XSsWs9NSxI
FjjE9/mOFVUH25KN+X7rzBPo2S0pMQLqyQxSLIdI2LPDxzlknctT6OoBPKPJjb7S
Lt5GhIA5Cz+cohfX6LePG4FkvwU32tTRBz5YNhFBizmS+YifBBgBAgAJBQJSx22G
AhsMAAoJEFVKIId/mLvOulED/2uUh/qjOT468XoK6Xt837w45JQPpLqiGH9KJgqF
rUxJMw1bIE2G606OY6hCgeE+YC8qny29hQtXhKIquUI/0A1qK3aCZhwqyqT+QjvF
6Xi0i/HrgQwCyBopY3uGndMbvthxU0KO0d6seMZltHDr8YaU1JvDwNFDQVuw+Rqy
57ET
=nvLl
-----END PGP PRIVATE KEY BLOCK-----`

const testKeys1And2Hex = "988d044d3c5c10010400b1d13382944bd5aba23a4312968b5095d14f947f600eb478e14a6fcb16b0e0cac764884909c020bc495cfcc39a935387c661507bdb236a0612fb582cac3af9b29cc2c8c70090616c41b662f4da4c1201e195472eb7f4ae1ccbcbf9940fe21d985e379a5563dde5b9a23d35f1cfaa5790da3b79db26f23695107bfaca8e7b5bcd0011010001b41054657374204b6579203120285253412988b804130102002205024d3c5c10021b03060b090807030206150802090a0b0416020301021e01021780000a0910a34d7e18c20c31bbb5b304009cc45fe610b641a2c146331be94dade0a396e73ca725e1b25c21708d9cab46ecca5ccebc23055879df8f99eea39b377962a400f2ebdc36a7c99c333d74aeba346315137c3ff9d0a09b0273299090343048afb8107cf94cbd1400e3026f0ccac7ecebbc4d78588eb3e478fe2754d3ca664bcf3eac96ca4a6b0c8d7df5102f60f6b0020003b88d044d3c5c10010400b201df61d67487301f11879d514f4248ade90c8f68c7af1284c161098de4c28c2850f1ec7b8e30f959793e571542ffc6532189409cb51c3d30dad78c4ad5165eda18b20d9826d8707d0f742e2ab492103a85bbd9ddf4f5720f6de7064feb0d39ee002219765bb07bcfb8b877f47abe270ddeda4f676108cecb6b9bb2ad484a4f0011010001889f04180102000905024d3c5c10021b0c000a0910a34d7e18c20c31bb1a03040085c8d62e16d05dc4e9dad64953c8a2eed8b6c12f92b1575eeaa6dcf7be9473dd5b24b37b6dffbb4e7c99ed1bd3cb11634be19b3e6e207bed7505c7ca111ccf47cb323bf1f8851eb6360e8034cbff8dd149993c959de89f8f77f38e7e98b8e3076323aa719328e2b408db5ec0d03936efd57422ba04f925cdc7b4c1af7590e40ab0020003988d044d3c5c33010400b488c3e5f83f4d561f317817538d9d0397981e9aef1321ca68ebfae1cf8b7d388e19f4b5a24a82e2fbbf1c6c26557a6c5845307a03d815756f564ac7325b02bc83e87d5480a8fae848f07cb891f2d51ce7df83dcafdc12324517c86d472cc0ee10d47a68fd1d9ae49a6c19bbd36d82af597a0d88cc9c49de9df4e696fc1f0b5d0011010001b42754657374204b6579203220285253412c20656e637279707465642070726976617465206b65792988b804130102002205024d3c5c33021b03060b090807030206150802090a0b0416020301021e01021780000a0910d4984f961e35246b98940400908a73b6a6169f700434f076c6c79015a49bee37130eaf23aaa3cfa9ce60bfe4acaa7bc95f1146ada5867e0079babb38804891f4f0b8ebca57a86b249dee786161a755b7a342e68ccf3f78ed6440a93a6626beb9a37aa66afcd4f888790cb4bb46d94a4ae3eb3d7d3e6b00f6bfec940303e89ec5b32a1eaaacce66497d539328b0020003b88d044d3c5c33010400a4e913f9442abcc7f1804ccab27d2f787ffa592077ca935a8bb23165bd8d57576acac647cc596b2c3f814518cc8c82953c7a4478f32e0cf645630a5ba38d9618ef2bc3add69d459ae3dece5cab778938d988239f8c5ae437807075e06c828019959c644ff05ef6a5a1dab72227c98e3a040b0cf219026640698d7a13d8538a570011010001889f04180102000905024d3c5c33021b0c000a0910d4984f961e35246b26c703ff7ee29ef53bc1ae1ead533c408fa136db508434e233d6e62be621e031e5940bbd4c08142aed0f82217e7c3e1ec8de574bc06ccf3c36633be41ad78a9eacd209f861cae7b064100758545cc9dd83db71806dc1cfd5fb9ae5c7474bba0c19c44034ae61bae5eca379383339dece94ff56ff7aa44a582f3e5c38f45763af577c0934b0020003"

const testKeys1And2PrivateHex = "9501d8044d3c5c10010400b1d13382944bd5aba23a4312968b5095d14f947f600eb478e14a6fcb16b0e0cac764884909c020bc495cfcc39a935387c661507bdb236a0612fb582cac3af9b29cc2c8c70090616c41b662f4da4c1201e195472eb7f4ae1ccbcbf9940fe21d985e379a5563dde5b9a23d35f1cfaa5790da3b79db26f23695107bfaca8e7b5bcd00110100010003ff4d91393b9a8e3430b14d6209df42f98dc927425b881f1209f319220841273a802a97c7bdb8b3a7740b3ab5866c4d1d308ad0d3a79bd1e883aacf1ac92dfe720285d10d08752a7efe3c609b1d00f17f2805b217be53999a7da7e493bfc3e9618fd17018991b8128aea70a05dbce30e4fbe626aa45775fa255dd9177aabf4df7cf0200c1ded12566e4bc2bb590455e5becfb2e2c9796482270a943343a7835de41080582c2be3caf5981aa838140e97afa40ad652a0b544f83eb1833b0957dce26e47b0200eacd6046741e9ce2ec5beb6fb5e6335457844fb09477f83b050a96be7da043e17f3a9523567ed40e7a521f818813a8b8a72209f1442844843ccc7eb9805442570200bdafe0438d97ac36e773c7162028d65844c4d463e2420aa2228c6e50dc2743c3d6c72d0d782a5173fe7be2169c8a9f4ef8a7cf3e37165e8c61b89c346cdc6c1799d2b41054657374204b6579203120285253412988b804130102002205024d3c5c10021b03060b090807030206150802090a0b0416020301021e01021780000a0910a34d7e18c20c31bbb5b304009cc45fe610b641a2c146331be94dade0a396e73ca725e1b25c21708d9cab46ecca5ccebc23055879df8f99eea39b377962a400f2ebdc36a7c99c333d74aeba346315137c3ff9d0a09b0273299090343048afb8107cf94cbd1400e3026f0ccac7ecebbc4d78588eb3e478fe2754d3ca664bcf3eac96ca4a6b0c8d7df5102f60f6b00200009d01d8044d3c5c10010400b201df61d67487301f11879d514f4248ade90c8f68c7af1284c161098de4c28c2850f1ec7b8e30f959793e571542ffc6532189409cb51c3d30dad78c4ad5165eda18b20d9826d8707d0f742e2ab492103a85bbd9ddf4f5720f6de7064feb0d39ee002219765bb07bcfb8b877f47abe270ddeda4f676108cecb6b9bb2ad484a4f00110100010003fd17a7490c22a79c59281fb7b20f5e6553ec0c1637ae382e8adaea295f50241037f8997cf42c1ce26417e015091451b15424b2c59eb8d4161b0975630408e394d3b00f88d4b4e18e2cc85e8251d4753a27c639c83f5ad4a571c4f19d7cd460b9b73c25ade730c99df09637bd173d8e3e981ac64432078263bb6dc30d3e974150dd0200d0ee05be3d4604d2146fb0457f31ba17c057560785aa804e8ca5530a7cd81d3440d0f4ba6851efcfd3954b7e68908fc0ba47f7ac37bf559c6c168b70d3a7c8cd0200da1c677c4bce06a068070f2b3733b0a714e88d62aa3f9a26c6f5216d48d5c2b5624144f3807c0df30be66b3268eeeca4df1fbded58faf49fc95dc3c35f134f8b01fd1396b6c0fc1b6c4f0eb8f5e44b8eace1e6073e20d0b8bc5385f86f1cf3f050f66af789f3ef1fc107b7f4421e19e0349c730c68f0a226981f4e889054fdb4dc149e8e889f04180102000905024d3c5c10021b0c000a0910a34d7e18c20c31bb1a03040085c8d62e16d05dc4e9dad64953c8a2eed8b6c12f92b1575eeaa6dcf7be9473dd5b24b37b6dffbb4e7c99ed1bd3cb11634be19b3e6e207bed7505c7ca111ccf47cb323bf1f8851eb6360e8034cbff8dd149993c959de89f8f77f38e7e98b8e3076323aa719328e2b408db5ec0d03936efd57422ba04f925cdc7b4c1af7590e40ab00200009501fe044d3c5c33010400b488c3e5f83f4d561f317817538d9d0397981e9aef1321ca68ebfae1cf8b7d388e19f4b5a24a82e2fbbf1c6c26557a6c5845307a03d815756f564ac7325b02bc83e87d5480a8fae848f07cb891f2d51ce7df83dcafdc12324517c86d472cc0ee10d47a68fd1d9ae49a6c19bbd36d82af597a0d88cc9c49de9df4e696fc1f0b5d0011010001fe030302e9030f3c783e14856063f16938530e148bc57a7aa3f3e4f90df9dceccdc779bc0835e1ad3d006e4a8d7b36d08b8e0de5a0d947254ecfbd22037e6572b426bcfdc517796b224b0036ff90bc574b5509bede85512f2eefb520fb4b02aa523ba739bff424a6fe81c5041f253f8d757e69a503d3563a104d0d49e9e890b9d0c26f96b55b743883b472caa7050c4acfd4a21f875bdf1258d88bd61224d303dc9df77f743137d51e6d5246b88c406780528fd9a3e15bab5452e5b93970d9dcc79f48b38651b9f15bfbcf6da452837e9cc70683d1bdca94507870f743e4ad902005812488dd342f836e72869afd00ce1850eea4cfa53ce10e3608e13d3c149394ee3cbd0e23d018fcbcb6e2ec5a1a22972d1d462ca05355d0d290dd2751e550d5efb38c6c89686344df64852bf4ff86638708f644e8ec6bd4af9b50d8541cb91891a431326ab2e332faa7ae86cfb6e0540aa63160c1e5cdd5a4add518b303fff0a20117c6bc77f7cfbaf36b04c865c6c2b42754657374204b6579203220285253412c20656e637279707465642070726976617465206b65792988b804130102002205024d3c5c33021b03060b090807030206150802090a0b0416020301021e01021780000a0910d4984f961e35246b98940400908a73b6a6169f700434f076c6c79015a49bee37130eaf23aaa3cfa9ce60bfe4acaa7bc95f1146ada5867e0079babb38804891f4f0b8ebca57a86b249dee786161a755b7a342e68ccf3f78ed6440a93a6626beb9a37aa66afcd4f888790cb4bb46d94a4ae3eb3d7d3e6b00f6bfec940303e89ec5b32a1eaaacce66497d539328b00200009d01fe044d3c5c33010400a4e913f9442abcc7f1804ccab27d2f787ffa592077ca935a8bb23165bd8d57576acac647cc596b2c3f814518cc8c82953c7a4478f32e0cf645630a5ba38d9618ef2bc3add69d459ae3dece5cab778938d988239f8c5ae437807075e06c828019959c644ff05ef6a5a1dab72227c98e3a040b0cf219026640698d7a13d8538a570011010001fe030302e9030f3c783e148560f936097339ae381d63116efcf802ff8b1c9360767db5219cc987375702a4123fd8657d3e22700f23f95020d1b261eda5257e9a72f9a918e8ef22dd5b3323ae03bbc1923dd224db988cadc16acc04b120a9f8b7e84da9716c53e0334d7b66586ddb9014df604b41be1e960dcfcbc96f4ed150a1a0dd070b9eb14276b9b6be413a769a75b519a53d3ecc0c220e85cd91ca354d57e7344517e64b43b6e29823cbd87eae26e2b2e78e6dedfbb76e3e9f77bcb844f9a8932eb3db2c3f9e44316e6f5d60e9e2a56e46b72abe6b06dc9a31cc63f10023d1f5e12d2a3ee93b675c96f504af0001220991c88db759e231b3320dcedf814dcf723fd9857e3d72d66a0f2af26950b915abdf56c1596f46a325bf17ad4810d3535fb02a259b247ac3dbd4cc3ecf9c51b6c07cebb009c1506fba0a89321ec8683e3fd009a6e551d50243e2d5092fefb3321083a4bad91320dc624bd6b5dddf93553e3d53924c05bfebec1fb4bd47e89a1a889f04180102000905024d3c5c33021b0c000a0910d4984f961e35246b26c703ff7ee29ef53bc1ae1ead533c408fa136db508434e233d6e62be621e031e5940bbd4c08142aed0f82217e7c3e1ec8de574bc06ccf3c36633be41ad78a9eacd209f861cae7b064100758545cc9dd83db71806dc1cfd5fb9ae5c7474bba0c19c44034ae61bae5eca379383339dece94ff56ff7aa44a582f3e5c38f45763af577c0934b0020000"

const dsaElGamalTestKeysHex = "9501e1044dfcb16a110400aa3e5c1a1f43dd28c2ffae8abf5cfce555ee874134d8ba0a0f7b868ce2214beddc74e5e1e21ded354a95d18acdaf69e5e342371a71fbb9093162e0c5f3427de413a7f2c157d83f5cd2f9d791256dc4f6f0e13f13c3302af27f2384075ab3021dff7a050e14854bbde0a1094174855fc02f0bae8e00a340d94a1f22b32e48485700a0cec672ac21258fb95f61de2ce1af74b2c4fa3e6703ff698edc9be22c02ae4d916e4fa223f819d46582c0516235848a77b577ea49018dcd5e9e15cff9dbb4663a1ae6dd7580fa40946d40c05f72814b0f88481207e6c0832c3bded4853ebba0a7e3bd8e8c66df33d5a537cd4acf946d1080e7a3dcea679cb2b11a72a33a2b6a9dc85f466ad2ddf4c3db6283fa645343286971e3dd700703fc0c4e290d45767f370831a90187e74e9972aae5bff488eeff7d620af0362bfb95c1a6c3413ab5d15a2e4139e5d07a54d72583914661ed6a87cce810be28a0aa8879a2dd39e52fb6fe800f4f181ac7e328f740cde3d09a05cecf9483e4cca4253e60d4429ffd679d9996a520012aad119878c941e3cf151459873bdfc2a9563472fe0303027a728f9feb3b864260a1babe83925ce794710cfd642ee4ae0e5b9d74cee49e9c67b6cd0ea5dfbb582132195a121356a1513e1bca73e5b80c58c7ccb4164453412f456c47616d616c2054657374204b65792031886204131102002205024dfcb16a021b03060b090807030206150802090a0b0416020301021e01021780000a091033af447ccd759b09fadd00a0b8fd6f5a790bad7e9f2dbb7632046dc4493588db009c087c6a9ba9f7f49fab221587a74788c00db4889ab00200009d0157044dfcb16a1004008dec3f9291205255ccff8c532318133a6840739dd68b03ba942676f9038612071447bf07d00d559c5c0875724ea16a4c774f80d8338b55fca691a0522e530e604215b467bbc9ccfd483a1da99d7bc2648b4318fdbd27766fc8bfad3fddb37c62b8ae7ccfe9577e9b8d1e77c1d417ed2c2ef02d52f4da11600d85d3229607943700030503ff506c94c87c8cab778e963b76cf63770f0a79bf48fb49d3b4e52234620fc9f7657f9f8d56c96a2b7c7826ae6b57ebb2221a3fe154b03b6637cea7e6d98e3e45d87cf8dc432f723d3d71f89c5192ac8d7290684d2c25ce55846a80c9a7823f6acd9bb29fa6cd71f20bc90eccfca20451d0c976e460e672b000df49466408d527affe0303027a728f9feb3b864260abd761730327bca2aaa4ea0525c175e92bf240682a0e83b226f97ecb2e935b62c9a133858ce31b271fa8eb41f6a1b3cd72a63025ce1a75ee4180dcc284884904181102000905024dfcb16a021b0c000a091033af447ccd759b09dd0b009e3c3e7296092c81bee5a19929462caaf2fff3ae26009e218c437a2340e7ea628149af1ec98ec091a43992b00200009501e1044dfcb1be1104009f61faa61aa43df75d128cbe53de528c4aec49ce9360c992e70c77072ad5623de0a3a6212771b66b39a30dad6781799e92608316900518ec01184a85d872365b7d2ba4bacfb5882ea3c2473d3750dc6178cc1cf82147fb58caa28b28e9f12f6d1efcb0534abed644156c91cca4ab78834268495160b2400bc422beb37d237c2300a0cac94911b6d493bda1e1fbc6feeca7cb7421d34b03fe22cec6ccb39675bb7b94a335c2b7be888fd3906a1125f33301d8aa6ec6ee6878f46f73961c8d57a3e9544d8ef2a2cbfd4d52da665b1266928cfe4cb347a58c412815f3b2d2369dec04b41ac9a71cc9547426d5ab941cccf3b18575637ccfb42df1a802df3cfe0a999f9e7109331170e3a221991bf868543960f8c816c28097e503fe319db10fb98049f3a57d7c80c420da66d56f3644371631fad3f0ff4040a19a4fedc2d07727a1b27576f75a4d28c47d8246f27071e12d7a8de62aad216ddbae6aa02efd6b8a3e2818cda48526549791ab277e447b3a36c57cefe9b592f5eab73959743fcc8e83cbefec03a329b55018b53eec196765ae40ef9e20521a603c551efe0303020950d53a146bf9c66034d00c23130cce95576a2ff78016ca471276e8227fb30b1ffbd92e61804fb0c3eff9e30b1a826ee8f3e4730b4d86273ca977b4164453412f456c47616d616c2054657374204b65792032886204131102002205024dfcb1be021b03060b090807030206150802090a0b0416020301021e01021780000a0910a86bf526325b21b22bd9009e34511620415c974750a20df5cb56b182f3b48e6600a0a9466cb1a1305a84953445f77d461593f1d42bc1b00200009d0157044dfcb1be1004009565a951da1ee87119d600c077198f1c1bceb0f7aa54552489298e41ff788fa8f0d43a69871f0f6f77ebdfb14a4260cf9fbeb65d5844b4272a1904dd95136d06c3da745dc46327dd44a0f16f60135914368c8039a34033862261806bb2c5ce1152e2840254697872c85441ccb7321431d75a747a4bfb1d2c66362b51ce76311700030503fc0ea76601c196768070b7365a200e6ddb09307f262d5f39eec467b5f5784e22abdf1aa49226f59ab37cb49969d8f5230ea65caf56015abda62604544ed526c5c522bf92bed178a078789f6c807b6d34885688024a5bed9e9f8c58d11d4b82487b44c5f470c5606806a0443b79cadb45e0f897a561a53f724e5349b9267c75ca17fe0303020950d53a146bf9c660bc5f4ce8f072465e2d2466434320c1e712272fafc20e342fe7608101580fa1a1a367e60486a7cd1246b7ef5586cf5e10b32762b710a30144f12dd17dd4884904181102000905024dfcb1be021b0c000a0910a86bf526325b21b2904c00a0b2b66b4b39ccffda1d10f3ea8d58f827e30a8b8e009f4255b2d8112a184e40cde43a34e8655ca7809370b0020000"

func TestXorKeys(t *testing.T) {
	shaedK1 := sha256.Sum256([]byte(key1))
	shaedK2 := sha256.Sum256([]byte(key2))
	shaedK3 := sha256.Sum256([]byte(key3))

	keys := [][]byte{shaedK1[:], shaedK2[:], shaedK3[:]}
	_, err := XorKeys(keys, kRequiredMaxKeySize)
	if err != nil {
		t.Fatalf("Unable to xor keys that's strange: %s.\n", err.Error())
	}

	keys = [][]byte{shaedK1[:], shaedK2[:], []byte(key3)}
	_, err = XorKeys(keys, kRequiredMaxKeySize)
	if err == nil {
		t.Fatalf("In this case xor should return an error key3 is too short.\n")
	}
}

func generateAesKey() ([]byte, error) {
	shaedK1 := sha256.Sum256([]byte(key1))
	shaedK2 := sha256.Sum256([]byte(key2))
	shaedK3 := sha256.Sum256([]byte(key3))

	keys := [][]byte{shaedK1[:], shaedK2[:], shaedK3[:]}
	aesk, err := XorKeys(keys, kRequiredMaxKeySize)
	if err != nil {
		return nil, fmt.Errorf("Unable to xor keys that's strange: %s.\n", err.Error())
	}
	return aesk, nil
}

func TestAESCBCEncryption(t *testing.T) {
	plainbytes := []byte(plaintex)

	// generate keys
	aesk, err := generateAesKey()
	if err != nil {
		t.Fatalf("%s", err)
	}
	derived := DeriveKeyWithPbkdf2(aesk, []byte(salt), 10000)

	// encrypt plain text
	cipheredbytes, err := AesEncrypt(derived, []byte(salt), plainbytes, CBC)
	if err != nil {
		t.Error("Unable to encrypt plain text:", err)
		return
	}
	t.Log("Encryption successfull.")

	// decrypt it
	decryptedbytes, err := AesDecrypt(derived, cipheredbytes, CBC)
	if err != nil {
		t.Error("Unable to decrypt ciphered bytes:", err)
		return
	}
	t.Log("Decryption successfull.")

	if string(decryptedbytes) != plaintex {
		t.Error("Encryption and decryption operations went wrong.")
		return
	}
	t.Log("Correctly AES CBC operations.")
}

func readerFromHex(s string) ([]byte, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func TestOpenPgpEncryption(t *testing.T) {
	data, _ := readerFromHex(testKeys1And2PrivateHex)
	keyr, err := openpgp.ReadKeyRing(bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Unable to extract private key: %d.\n", err.Error())
	}

	plainbytes := []byte(plaintex)
	signer := keyr[0]
	encrypted, err := OpenPgpEncrypt(plainbytes, keyr[:1], signer)
	if err != nil {
		t.Fatalf("Unexpected error: %s.\n", err.Error())
	}

	encoded, _ := EncodePgpArmored(encrypted, kEn1gm4Type)
	t.Logf("Encrypted blob: %s.\n", string(encoded))

	plaintdata, err := OpenPgpDecrypt(encrypted, keyr)
	if err != nil {
		t.Fatalf("Unable to access body: %s.\n", err.Error())
	}
	if bytes.Compare(plaintdata, plainbytes) != 0 {
		t.Fatalf("Unexpected result are different\n")
	}
	t.Logf("Plaintext: %s.\n", string(plaintdata))
}

func TestOpenPgpArmoredKeys(t *testing.T) {
	entityList, err := ReadArmoredKeyRing([]byte(publicKey), nil)
	if err != nil {
		t.Fatalf("Unable to access public armored key: %s.\n", err.Error())
	}
	if len(entityList) != 1 {
		t.Fatalf("Unexpected size having %d expecting %d.\n", len(entityList), 1)
	}
	pvkList, err := ReadArmoredKeyRing([]byte(privateKey), []byte("golang"))
	if err != nil {
		t.Fatalf("Unable to access private armored key: %s.\n", err.Error())
	}
	if len(pvkList) != 1 {
		t.Fatalf("Unexpected size having %d expecting %d.\n", len(pvkList), 1)
	}
	entityList = append(entityList, pvkList...)
	plainbytes := []byte(plaintex)
	_, err = OpenPgpEncrypt(plainbytes, entityList, pvkList[0])
	if err != nil {
		t.Fatalf("Unexpected error: %s.\n", err.Error())
	}
}

func TestOpenPgpSignature(t *testing.T) {
	pvkList, err := ReadArmoredKeyRing([]byte(privateKey), []byte("golang"))
	if err != nil {
		t.Fatalf("Unable to access private armored key: %s.\n", err.Error())
	}
	if len(pvkList) != 1 {
		t.Fatalf("Unexpected size having %d expecting %d.\n", len(pvkList), 1)
	}
	msg := []byte(plaintex)
	signature, err := OpenPgpSignMessage(msg, pvkList[0])
	if err != nil {
		t.Fatalf("Unable to sign the message: %s.\n", err.Error())
	}
	if len(signature) == 0 {
		t.Fatalf("Invalid signature len should be not 0.\n")
	}
	t.Logf("Signature: %v (%d).\n", signature, len(signature))
	// test valid signer key
	err = OpenPgpVerifySignature(signature, msg, pvkList[0])
	if err != nil {
		t.Fatalf("Unable to verify signature: %s.\n", err.Error())
	}
	// test other user's keys
	entityList, err := ReadArmoredKeyRing([]byte(publicKey), nil)
	if err != nil {
		t.Fatalf("Unable to access public armored key: %s.\n", err.Error())
	}
	if len(entityList) == 0 {
		t.Fatalf("Expected at least a public key.\n")
	}
	err = OpenPgpVerifySignature(signature, msg, entityList[0])
	if err == nil {
		t.Fatalf("This user should not be able to verify the signature.\n")
	}
}

func TestKeysCreation(t *testing.T) {
	pvk, pbk, err := NewPgpKeypair("user", "This is a test key", "user@mail.com")
	if err != nil {
		t.Fatalf("Unable to create keys: %s.\n", err.Error())
	}
	t.Logf("Keys: %s pub %s.\n", hex.EncodeToString(pvk), hex.EncodeToString(pbk))

	pk, err := openpgp.ReadKeyRing(bytes.NewBuffer(pvk))
	if err != nil {
		t.Fatalf("Unable to extract private key: %s.\n", err.Error())
	}
	if len(pk) != 1 {
		t.Fatalf("Unexpected single private key not found.\n")
	}
	pb, err := openpgp.ReadKeyRing(bytes.NewBuffer(pbk))
	if err != nil {
		t.Fatalf("Unable to extract public key: %d.\n", err.Error())
	}
	if len(pb) != 1 {
		t.Fatalf("Unexpected single private key not found.\n")
	}

	plainbytes := []byte(plaintex)

	encrypted, err := OpenPgpEncrypt(plainbytes, pb, pk[0])
	if err != nil {
		t.Fatalf("Unexpected error: %s.\n", err.Error())
	}

	encoded, _ := EncodePgpArmored(encrypted, kEn1gm4Type)
	t.Logf("Encrypted blob: %s.\n", string(encoded))

	plaintdata, err := OpenPgpDecrypt(encrypted, pk)
	if err != nil {
		t.Fatalf("Unable to access body: %s.\n", err.Error())
	}
	if bytes.Compare(plaintdata, plainbytes) != 0 {
		t.Fatalf("Unexpected result are different\n")
	}
	t.Logf("Plaintext: %s.\n", string(plaintdata))
}
