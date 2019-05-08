package main

import (
	"os"
	"testing"

	lambdat "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/lambda-testing"
)

func TestUser(t *testing.T) {
	user := os.Getenv("USER")
	t.Log(user)
}

func TestHandler(t *testing.T) {
	//this line is for the token encryption/decription should be present when function is used to authentiacate
	os.Setenv("RSAPRIVATEKEY", "308204a40201000282010100d6a54d6b1209cba6712c73bbbaf5e6fc3678989c5bffde5e7b8de825018d65dfa43a84dc9e8e021d450c4e65f1a2ad3674e0c2c7572ca8895051b5931f176ae03c9d42c7ccc87071c40ff39daecd632a53db2297d894edc8e2c5392be7b571149222648fd096bb1c01be92ce64009c7da2ae090599f370090852c2436f05679f9433ffa99dee3a3252d81084e14f7613b2775ff20bb0cbb05cd17713991ee74e05fd0770ed45354d27ef6c93688f8b478cba4c4ffef75bc5af1820e2f8b175cd47cc1ac4dd4095a6f75f9dca4ad862ac762ebdd08a6173690e91d76a9f704bba4926d9ce884f0ccd83ab77550bf32414a94e403582db55d8b9b7eb741f525a45020301000102820101009c83b4cc1a3f224c9fc1b63271c5d5449bc39c2487c12fb8dd87407b9b822b82c41217c777a63d4c7288e2b1db5cafc941b892cf2075e3ff1c9e3834ab3b3c277e8b7da28b64acf987e9c9ce752436e72a7663e72d7a8b592c627ba9d42fade13e1dee0e201f891886fd1bb77b9c2680461b7960a83da6b82f658959fa9e8a4bafb12b65306842d0c623eb37e9cc29a35a19fd237b6f0d2bb2864c32e407788dab912e755950724ecfd4d453770eb7f451ad31f998fb49d838e4e65a7753e09c883e91f4ce8fce5fff56aecb574d7d6fdb1fecdd1afbbbdd3e677d5fe0f5816d7a4421984c799521fd461e1043603a072cbf30e5e5ed9a8b2d0ecaf2669db11102818100fd7d3ff870e1670a6328a2e17e0dfc40502a1a8cd6f2a3da27e2c271cf0cf6cc8d4ae0af7f0a7a35701d1e0520c58cef4917894a80fea46bbec6971e675eca24dfaa1c406c95c8de07da7356d9696b1df49c9c2071bc88ad7fff5aee55e8ea1a770bd53939356d9d48e12621314303f59af6a3504b1da2b5a4841bd6553859ab02818100d8c58f776dbbbedbb0bf62849c98bbce17b9ade44555fdb1d2de0a9db7803467a66de9f36c5b065a06b326557a69614de7df8310bf52767dd0e0ae594a4c4e87989097f765fe26555d89d8dd22d4c9ede684e91ccd1b29a3d2f1e7fe36c59de81952d12cdba2d684892c5a8dc79d4ddf8cd1a2e92ab1ad1715d5ad8a1c368bcf02818100c5758a9e51f813571114f78465b8293643fbf8409bb3d915381ab8d304b199928fc1b332a1e89c7802047c7d0c2136feb2d625b926b0b58dc4c757b2745d6f63b7e3002ce328ee969e5179a53ea892ab7bff7ed2fb261ce5e21e1d4c2919cd3a9e5f565244112d78e6eb93d3295785bf0d5e70ab3c483296023872a2cc31a00f02818021e3d12834c9b36f1954f28c150773e526a46ae1534dbc59fec3a419404514ec5782bb9ec903fa1c3d0be92457fcfdaf765ee558caf09381dc14246de545c4c9423ae8e74ed4cb1d7180499d5902b7873010fb78fb4011e480e83d02eb813dccb998cf071a577cfe3f8be5a460dee0fbe0422e1c1206b12ef8c4ed5ab84a76d50281804f7cdd2c3c5e25472b6c9b94a284ce5a64f5aa908051a765edfe49b8969d5342770cc2e1889ae4de8851a3412c6089b48cb3c4e6925bfd8b63de205110af51eb579ebaf9be95885be8e63c899a48df5a3de284fc760c913f68ed297e2a83543ed3156acf59fdb6b91daf85ad61fabe298307e7e5eff4cf0ffed05728f6c78954")
	body := map[string]interface{}{
		// "token": "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.BnoC2HMe15LI35Yqpl_jk0cCn50y7MByanOKyoeYXyCQINTyauabHPHCgxOJlhh_2Uy_15jaxZf7C_bExixKCwq42XURTmqq1Kkf2dacaSf_ZULuH0c0sz7PMQHsaFMPsQz8StXqtqepOzEVy-Iq-Y-aFq_viAogQ9r7y6OJjjghpDgRx9hmcT_ji7slAZgTOoACAhCz4oV5fEjTmdQXgNZNvXalEyL2WVe4ArKj6mt2IcdYnDqdhO5kj27M10YY90kOqUUWSTTkBXIfg08U-REchaHBVqLdK65LMZd2ZQfKK5TVJlndNC8pDCXl99L9qogC15dGcVn9r5S4FF67ug.qJjEh6nFsVv3gJ5E.6mp7WdsOhRvqMr5HYFeO5Thryiy0fquoL0UGHyVorsbOyObyXL5b87H6gtaNB0mu7g.qDfaSrZEuF0SPEY1SZ3gtQ",
		// "token": "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.q6KtzF77TckTDlXHMMhHfP5EEqscEN4DS9zLyXHfcbvxfIE3oeGQ7gxBGrrNPxN4Zl_AkMhDT6D6jP0xkAaa8zjTUR9iEtYfS2TVlt33Ei-8m5WG7pNp4NbyIX1-nEwO-GhQKQkcP5dH-SLsEHioDsxUyp4FFBN70ge22Ryoemk60jk9IZw-bt1VrnI1pjvInLi_k0_jxqGlZ3FA7Q1-oyIs7LUloG3GQtmpLb0Pv0gCz6yfXILq3XAEdNoPyZo6YFn6lmTN7nAHXjCnkH5HBf-cOEjW0YymgIQ3qzApWEYQv19v25jCHBNztpw2LvlPJ_eFNQrXy3jhrZqyvkFH2Q.ebPXqBQ4QJN0MAxp.ZLId7IUne2vPKklnqJx5LGL0AGMx69TEaf4_jnHwMRiAokBG-CEFw06rLhYXgu3-.YV4qMlHjZ3MlyfemKvFVZQ",
		// "token": "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.0dpbnRw9A8q_cgr7oWArTByuv4KOBDE-jL1M7_SyuGmcrT9RN4bvIKW9s96ji8DlCotMArrLWIwOzU5-eQHP8Hlap6TdIMDU7s1V1m_THWW9rXXJ-EVFMQotgKXwAjkrrjxC8awfvXDzstrV1pFPYH6mXB-FRPz9rKtkwmZJVrVq2pUVKTxBq1DDUkUJP0NOMyWiz2bejkVOtjaGDTZUK4PfgXmNNV-BShHmxEgJH5_UHNl1pFSJJEy87wQ5p3rjdZb-f7ZAqnx_VEtexVkeXMRIc6UG6mgM-_3GAqg7g2jAJLpqDjLppDTyRTdhygihdCFURQ2-9n317nZ5y2iyXg.6IvFeYsGnMGSOO5W.SFFcPkU3bhW2xxkwtCrX9R0NBJNCF_YP1oC-obzmP_geMqgo6QnSAKST5oQLrICAD1A.jmbXUkUpxu4t64P2ywI81g",
		"token": "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.c61wP5tNKMwGJYlKM1fCvF9q5TNza_e-54O5DhwMS_fouSOINFuIcsFniiCbSH8zns60IbwaLZuOoww0PDAdPuGny8lmD3I6T4apV7afnZVq2DfSXxwYugr2TUiKtuQLgXBob6Sk9xfQBLrrBTxzmtywJejcIuWJWYYNSZfR2DjY9apUXy9xi4uOsPLWg_8QejbROMhG_dmaAotMzl8x8WX3lb09AOnXCX6vJxx-INOPCarY581R2kAAtBVHqlDDDtX1YPgL8vdIHlgK7sFf2kF4juNYqRCHvM6FfKdT8V_RlvMWcqfN8lyxsy0nA08oHxhpmTEHIT4jt4T5DbyuTA.YgiMp2kzA_l2Oh-f.YBCgYTua2FEZevrQj6vCjazOgZvV2S60fRY.K5GjAfxd3Q7fBMhanHU36A",

		"group_id": "bfdcc08eab2351be",
		"member":   "kachamaka4",
	}
	req, err := lambdat.ReadBody(body)
	if err != nil {
		t.Error(err)
	}
	res, err := lambdat.InvokeHandler(handler, req)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
