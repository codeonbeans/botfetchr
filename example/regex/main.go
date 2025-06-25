package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func main() {
	// text := `t":0.000000,"url144":"https:\/\/vk6-7.vkuser.net\/?srcIp=58.187.246.12&pr=40&expires=1751116459527&srcAg=CHROME&fromCache=1&ms=95.142.206.166&type=4&subId=8151226255876&sig=UODdnK5TZkU&ct=0&urls=45.136.22.165%3B185.226.52.209&clientType=13&appId=512000384397&zs=43&id=8150980758020"`

	// re := regexp.MustCompile(`"url\d+":"([^"]+)"`)
	// match := re.FindStringSubmatch(text)
	// if len(match) > 1 {
	// 	fmt.Println("Found URL:", match[1])
	// } else {
	// 	fmt.Println("No match found")
	// }

	// text := `"https:\\/\\/vkvd296.okcdn.ru\\/?srcIp=58.187.246.12&pr=40&expires=1751129922986&srcAg=CHROME&fromCache=1&ms=185.226.53.133&type=5&subId=7268260383481&sig=RSu0I7SJ5BE&ct=0&urls=45.136.21.154&clientType=13&appId=512000384397&zs=65&id=7266165197561"`
	// var idk string
	// json.Unmarshal([]byte(text), &idk)
	// idk, _ = FixEscapedURL(idk)
	// fmt.Println("NIGGA KYS", idk)

	// resp, err := http.Get(idk)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	text := `"result": {
                                                    "data": {
                                                        "xdt_api__v1__media__shortcode__web_info": {
                                                            "items": [
                                                                {
                                                                    "code": "DLAi8xlySmh",
                                                                    "pk": "3657076607422114209",
                                                                    "id": "3657076607422114209_14041886518",
                                                                    "ad_id": null,
                                                                    "taken_at": 1750177529,
                                                                    "inventory_source": null,
                                                                    "video_versions": [
                                                                        {
                                                                            "width": 480,
                                                                            "height": 853,
                                                                            "url": "https:\/\/instagram.fsgn2-9.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQOJ81KLmW5UxKmEtz_GK8wTf0FvEkzuX36Qf8GoAI_BZxp65gXFz20WScdaQFNoAHbMIwEwv1bn_9XByiXhodMW.mp4?stp=dst-mp4&efg=eyJ2aWRlb19pZCI6IjEyMzM2MDk5MjgyNTU5NzciLCJxZV9ncm91cHMiOiJbXCJpZ19wcm9ncmVzc2l2ZV91cmxnZW4ucHJvZHVjdF90eXBlLmNsaXBzXCJdIiwidmVuY29kZV90YWciOiJpZ19wcm9ncmVzc2l2ZV91cmxnZW4uY2xpcHMuYzIuc3ZlX3NkIn0&_nc_cat=106&ccb=9-4&oh=00_AfN3BgBdO5YMr-wDwTF1LJwaykf-U3ViZyyzQB8jFh5OBw&oe=685CB59D&_nc_sid=10d13b",
                                                                            "type": 101
                                                                        },
                                                                        {
                                                                            "width": 480,
                                                                            "height": 853,
                                                                            "url": "https:\/\/instagram.fsgn2-9.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQOJ81KLmW5UxKmEtz_GK8wTf0FvEkzuX36Qf8GoAI_BZxp65gXFz20WScdaQFNoAHbMIwEwv1bn_9XByiXhodMW.mp4?stp=dst-mp4&efg=eyJ2aWRlb19pZCI6IjEyMzM2MDk5MjgyNTU5NzciLCJxZV9ncm91cHMiOiJbXCJpZ19wcm9ncmVzc2l2ZV91cmxnZW4ucHJvZHVjdF90eXBlLmNsaXBzXCJdIiwidmVuY29kZV90YWciOiJpZ19wcm9ncmVzc2l2ZV91cmxnZW4uY2xpcHMuYzIuc3ZlX3NkIn0&_nc_cat=106&ccb=9-4&oh=00_AfN3BgBdO5YMr-wDwTF1LJwaykf-U3ViZyyzQB8jFh5OBw&oe=685CB59D&_nc_sid=10d13b",
                                                                            "type": 103
                                                                        },
                                                                        {
                                                                            "width": 480,
                                                                            "height": 853,
                                                                            "url": "https:\/\/instagram.fsgn2-9.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQOJ81KLmW5UxKmEtz_GK8wTf0FvEkzuX36Qf8GoAI_BZxp65gXFz20WScdaQFNoAHbMIwEwv1bn_9XByiXhodMW.mp4?stp=dst-mp4&efg=eyJ2aWRlb19pZCI6IjEyMzM2MDk5MjgyNTU5NzciLCJxZV9ncm91cHMiOiJbXCJpZ19wcm9ncmVzc2l2ZV91cmxnZW4ucHJvZHVjdF90eXBlLmNsaXBzXCJdIiwidmVuY29kZV90YWciOiJpZ19wcm9ncmVzc2l2ZV91cmxnZW4uY2xpcHMuYzIuc3ZlX3NkIn0&_nc_cat=106&ccb=9-4&oh=00_AfN3BgBdO5YMr-wDwTF1LJwaykf-U3ViZyyzQB8jFh5OBw&oe=685CB59D&_nc_sid=10d13b",
                                                                            "type": 102
                                                                        }
                                                                    ],
                                                                    "coauthor_producers": [
                                                                    ],
                                                                    "invited_coauthor_producers": [
                                                                    ],
                                                                    "facepile_top_likers": null,
                                                                    "is_dash_eligible": 1,
                                                                    "number_of_qualities": 2,
                                                                    "video_dash_manifest": "\u003C?xml version=\"1.0\" encoding=\"UTF-8\"?>\n\u003CMPD xmlns=\"urn:mpeg:dash:schema:mpd:2011\" xmlns:xsi=\"http:\/\/www.w3.org\/2001\/XMLSchema-instance\" xsi:schemaLocation=\"urn:mpeg:dash:schema:mpd:2011 DASH-MPD.xsd\" profiles=\"urn:mpeg:dash:profile:isoff-on-demand:2011\" minBufferTime=\"PT2S\" type=\"static\" mediaPresentationDuration=\"PT53.066666S\" FBManifestIdentifier=\"FgAYCmJhc2ljX2dlbjIZNpjMg7XPio8C2vzD7MSrowKkhInogsa4AiIYGGRhc2hfbG5faGVhYWNfdmJyM19hdWRpbyIA\">\u003CPeriod id=\"0\" duration=\"PT53.066666S\">\u003CAdaptationSet id=\"0\" contentType=\"video\" subsegmentAlignment=\"true\" par=\"9:16\" FBUnifiedUploadResolutionMos=\"360:77\">\u003CSupplementalProperty schemeIdUri=\"urn:mpeg:mpegB:cicp:TransferCharacteristics\" value=\"1\"\/>\u003CRepresentation id=\"687298224136466vd\" bandwidth=\"322934\" codecs=\"avc1.4d001e\" mimeType=\"video\/mp4\" sar=\"1:1\" FBEncodingTag=\"dash_h264-basic-gen2_360p\" FBContentLength=\"2142132\" FBPlaybackResolutionMos=\"0:100,360:69.1,480:60.8,720:49.8,1080:37.9\" FBPlaybackResolutionMosConfidenceLevel=\"high\" FBPlaybackResolutionCsvqm=\"0:100,360:83.9,480:76.1,720:64.7,1080:51.7\" FBAbrPolicyTags=\"\" width=\"360\" height=\"640\" frameRate=\"15360\/512\" FBDefaultQuality=\"1\" FBQualityClass=\"sd\" FBQualityLabel=\"360p\">\u003CBaseURL>https:\/\/instagram.fsgn2-11.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQMT4rTYIrgulIDgdq2_CuTq0fHmn7f7Gx4uNBkAB8NtER0WLjJf4q89ItEBGBstPZAi3CBCTH-gPzuYKqVMZ7lM.mp4?strext=1&amp;_nc_cat=105&amp;_nc_oc=AdmDMx0jcuIYHCyflt_TpB6GTWq7C-0qC5tD9w4RCcoP2SN2tmzDmpx2364JBfBUj1g&amp;_nc_sid=9ca052&amp;_nc_ht=instagram.fsgn2-11.fna.fbcdn.net&amp;_nc_ohc=2ke0Hrhs1fEQ7kNvwF_1Oqb&amp;efg=eyJ2ZW5jb2RlX3RhZyI6ImlnLXhwdmRzLmNsaXBzLmMyLUMzLmRhc2hfaDI2NC1iYXNpYy1nZW4yXzM2MHAiLCJ2aWRlb19pZCI6bnVsbCwib2lsX3VybGdlbl9hcHBfaWQiOjEyMTc5ODE2NDQ4Nzk2MjgsImNsaWVudF9uYW1lIjoiaWciLCJ4cHZfYXNzZXRfaWQiOjYwNDQ5NjY3OTA4Nzg3MiwidmlfdXNlY2FzZV9pZCI6MTAxMjAsImR1cmF0aW9uX3MiOjUzLCJ1cmxnZW5fc291cmNlIjoid3d3In0\u00253D&amp;ccb=17-1&amp;_nc_zt=28&amp;oh=00_AfPppvHid4zsKru4hxMPKsu-uiiCRFOBwcTCQMpWgKjTEw&amp;oe=6860B517\u003C\/BaseURL>\u003CSegmentBase indexRange=\"887-1050\" timescale=\"15360\" FBMinimumPrefetchRange=\"1051-13717\" FBPartialPrefetchDuration=\"2500\" FBPartialPrefetchRange=\"1051-105662\" FBFirstSegmentRange=\"1051-192752\" FBFirstSegmentDuration=\"5000\" FBSecondSegmentRange=\"192753-415543\" FBPrefetchSegmentRange=\"1051-192752\" FBPrefetchSegmentDuration=\"5000\">\u003CInitialization range=\"0-886\"\/>\u003C\/SegmentBase>\u003C\/Representation>\u003CRepresentation id=\"640663742349101v\" bandwidth=\"1571482\" codecs=\"avc1.64001f\" mimeType=\"video\/mp4\" sar=\"1:1\" FBEncodingTag=\"dash_h264-basic-gen2_720p\" FBContentLength=\"10424170\" FBPlaybackResolutionMos=\"0:100,360:95.4,480:92.5,720:83.4,1080:73.9\" FBPlaybackResolutionMosConfidenceLevel=\"high\" FBPlaybackResolutionCsvqm=\"0:100,360:98.56,480:97.2,720:94.4,1080:88.4\" FBAbrPolicyTags=\"\" width=\"720\" height=\"1280\" frameRate=\"15360\/512\" FBQualityClass=\"hd\" FBQualityLabel=\"720p\">\u003CBaseURL>https:\/\/instagram.fsgn2-3.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQNaCLFa3BE9-BKjJ7o6kgImGsJ5tj3eqmPj4DLdxWos0ShK-hwHUdsk3wMZYt5Ve-NPUFWIyTcJ0dHcvrAynmCF.mp4?strext=1&amp;_nc_cat=107&amp;_nc_oc=Adl_6hDVkzsIEclzySp1rRzt1LgOVXJ5Dlz_mZSQ_Knawij9a4jLcPYAl3pHCCbST7c&amp;_nc_sid=9ca052&amp;_nc_ht=instagram.fsgn2-3.fna.fbcdn.net&amp;_nc_ohc=wbgDHM7ziYcQ7kNvwFyr_lO&amp;efg=eyJ2ZW5jb2RlX3RhZyI6ImlnLXhwdmRzLmNsaXBzLmMyLUMzLmRhc2hfaDI2NC1iYXNpYy1nZW4yXzcyMHAiLCJ2aWRlb19pZCI6bnVsbCwib2lsX3VybGdlbl9hcHBfaWQiOjEyMTc5ODE2NDQ4Nzk2MjgsImNsaWVudF9uYW1lIjoiaWciLCJ4cHZfYXNzZXRfaWQiOjYwNDQ5NjY3OTA4Nzg3MiwidmlfdXNlY2FzZV9pZCI6MTAxMjAsImR1cmF0aW9uX3MiOjUzLCJ1cmxnZW5fc291cmNlIjoid3d3In0\u00253D&amp;ccb=17-1&amp;_nc_zt=28&amp;oh=00_AfMTRoqQ-NjOidbm9tTxhgCLZNdjn99SV2_BO0wLMIykKg&amp;oe=6860C515\u003C\/BaseURL>\u003CSegmentBase indexRange=\"892-1055\" timescale=\"15360\" FBMinimumPrefetchRange=\"1056-34291\" FBPartialPrefetchDuration=\"2500\" FBPartialPrefetchRange=\"1056-500462\" FBFirstSegmentRange=\"1056-955466\" FBFirstSegmentDuration=\"5000\" FBSecondSegmentRange=\"955467-2016634\" FBPrefetchSegmentRange=\"1056-955466\" FBPrefetchSegmentDuration=\"5000\">\u003CInitialization range=\"0-891\"\/>\u003C\/SegmentBase>\u003C\/Representation>\u003C\/AdaptationSet>\u003CAdaptationSet id=\"1\" contentType=\"audio\" subsegmentStartsWithSAP=\"1\" subsegmentAlignment=\"true\">\u003CRepresentation id=\"596117759750924ad\" bandwidth=\"53184\" codecs=\"mp4a.40.5\" mimeType=\"audio\/mp4\" FBAvgBitrate=\"53184\" audioSamplingRate=\"44100\" FBEncodingTag=\"dash_ln_heaac_vbr3_audio\" FBContentLength=\"353891\" FBPaqMos=\"85.40\" FBAbrPolicyTags=\"\" FBDefaultQuality=\"1\">\u003CAudioChannelConfiguration schemeIdUri=\"urn:mpeg:dash:23003:3:audio_channel_configuration:2011\" value=\"2\"\/>\u003CBaseURL>https:\/\/instagram.fsgn2-3.fna.fbcdn.net\/o1\/v\/t16\/f2\/m69\/AQO7LAN_PgKgL2_mKt2rBazMDu2jm6-PWA_sHn9ekGO8GQOXO0i6lY64gVafTuo8fa7TkbbyWkjTUu4aFQaxa7zR.mp4?strext=1&amp;_nc_cat=107&amp;_nc_oc=AdnWWCFn1NQnZg5p-1pZATOcPGZbYGuY__tHgGgcn8MNfO7POHUzDRRtvod6ru227ds&amp;_nc_sid=9ca052&amp;_nc_ht=instagram.fsgn2-3.fna.fbcdn.net&amp;_nc_ohc=MCerWwiLODQQ7kNvwFbshB3&amp;efg=eyJ2ZW5jb2RlX3RhZyI6ImlnLXhwdmRzLmNsaXBzLmMyLUMzLmRhc2hfbG5faGVhYWNfdmJyM19hdWRpbyIsInZpZGVvX2lkIjpudWxsLCJvaWxfdXJsZ2VuX2FwcF9pZCI6MTIxNzk4MTY0NDg3OTYyOCwiY2xpZW50X25hbWUiOiJpZyIsInhwdl9hc3NldF9pZCI6NjA0NDk2Njc5MDg3ODcyLCJ2aV91c2VjYXNlX2lkIjoxMDEyMCwiZHVyYXRpb25fcyI6NTMsInVybGdlbl9zb3VyY2UiOiJ3d3cifQ\u00253D\u00253D&amp;ccb=17-1&amp;_nc_zt=28&amp;oh=00_AfO-SIuKf3aIHSRzGzu0arVONmbkhf-vsRXyVzmKBd2fEA&amp;oe=6860A2CD\u003C\/BaseURL>\u003CSegmentBase indexRange=\"824-1179\" timescale=\"44100\" FBMinimumPrefetchRange=\"1180-1523\" FBPartialPrefetchDuration=\"2500\" FBPartialPrefetchRange=\"1180-18177\" FBFirstSegmentRange=\"1180-15567\" FBFirstSegmentDuration=\"2021\" FBSecondSegmentRange=\"15568-29466\" FBPrefetchSegmentRange=\"1180-29466\" FBPrefetchSegmentDuration=\"4017\">\u003CInitialization range=\"0-823\"\/>\u003C\/SegmentBase>\u003C\/Representation>\u003C\/AdaptationSet>\u003C\/Period>\u003C\/MPD>\n",
                                                                    "image_versions2": {
                                                                        "candidates": [
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e15_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfPWrwIas8ipqZ_FEQRPptd_jiX8dBEA12rtcYrUSJJekw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 1920,
                                                                                "width": 1080
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e35_p720x720_sh0.08_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfOEbWaF3Od6YsGvMA8K0Unsp90NdRA0rlYCmef2sJZSBQ&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 1280,
                                                                                "width": 720
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e35_p640x640_sh0.08_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfMvYqS0rskgA9KFfITjh6ZoKFNZAVyqCUn2PFiBHH9JRQ&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 1138,
                                                                                "width": 640
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e15_p480x480_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfNjlTbMh_zxXX-6BawrYioXSAE5DF8Zf8EOwQmIIrfpLw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 853,
                                                                                "width": 480
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e15_p320x320_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfO134tvYvtZzpgPmWoyhAFteknIzol_qnxws_QlufL97g&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 569,
                                                                                "width": 320
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=dst-jpg_e15_p240x240_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfMwT75pRVxspOubQrjK601gKTilgOR6tni77_y1_7mjLw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 427,
                                                                                "width": 240
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e15_fr_s1080x1080_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfNOLZ-YB3hryN6PGfu7QzIvRHhYnGMD1VWXrFD92xuokw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 1080,
                                                                                "width": 1080
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e35_s750x750_sh0.08_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfMOrslRYzyLNOjM__uodU-rmKU4dbi9El2Yqvp1gUpDvw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 750,
                                                                                "width": 750
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e35_s640x640_sh0.08_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfO2EhaXl37Z_vcAbF4xLWOHVv0WOuTA0jG5tAFN04Eksg&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 640,
                                                                                "width": 640
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e15_s480x480_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfO3LJH_EkjMIGo0lwbRS2YXK1enybK8hslXfU5vA75eYA&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 480,
                                                                                "width": 480
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e15_s320x320_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfNfbWV56JNksCXiESz5VrLqpNBkrcEGvmaXyAYZ4LmSzw&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 320,
                                                                                "width": 320
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e15_s240x240_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfODJNwzxnCSQFa7vNniAkzmfKmQzIpB30x2yzMK2mmtgA&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 240,
                                                                                "width": 240
                                                                            },
                                                                            {
                                                                                "url": "https:\/\/instagram.fsgn2-8.fna.fbcdn.net\/v\/t15.5256-10\/509298119_712860918029462_5626852161443132666_n.jpg?stp=c0.420.1080.1080a_dst-jpg_e15_s150x150_tt6&efg=eyJ2ZW5jb2RlX3RhZyI6IkNMSVBTLmltYWdlX3VybGdlbi4xMDgweDE5MjAuc2RyLmY1MjU2LmZpcnN0X2ZyYW1lX3RodW1ibmFpbCJ9&_nc_ht=instagram.fsgn2-8.fna.fbcdn.net&_nc_cat=102&_nc_oc=Q6cZ2QGfBcQry_o3d8JwGEfxYZWo-shQe66kgF61r6bfkGI1fQTVueKAyXUIRi2HYJZrdKY&_nc_ohc=G1Qz_5YwFd8Q7kNvwEuIt7E&_nc_gid=QRVlt8lIXUV3V8FUkkD3mA&edm=APs17CUAAAAA&ccb=7-5&ig_cache_key=MzY1NzA3NjYwNzQyMjExNDIwOQ\u00253D\u00253D.3-ccb7-5&oh=00_AfOXGSg6z60OBvAMVHjnUDWAVsi-AlTt9PrdlgWlQQPxxQ&oe=6860C02C&_nc_sid=10d13b",
                                                                                "height": 150,
                                                                                "width": 150
                                                                            }
                                                                        ]
                                                                    },
                                                                    "is_paid_partnership": false,
                                                                    "sponsor_tags": null,
                                                                    "original_height": 1920,`

	matches := extractVideoURLs(text)
	fmt.Println("Extracted URLs:", strings.Join(matches, "\n\n"))
}

func UnmarshalURL(marshalledURL string) (string, error) {
	var result string
	err := json.Unmarshal([]byte(marshalledURL), &result)
	return result, err
}

func FixEscapedURL(escapedURL string) (string, error) {
	// Replace escaped forward slashes
	fixedURL := strings.ReplaceAll(escapedURL, `\/`, `/`)

	// Validate the URL
	_, err := url.Parse(fixedURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL after fixing: %v", err)
	}

	return fixedURL, nil
}

func extractVideoURLs(jsonData string) []string {
	// Regex pattern to match URLs within video_versions array
	pattern := `(?s)"video_versions":\s*\[(.*?)\]`
	videoVersionsRegex := regexp.MustCompile(pattern)

	// Find the video_versions section
	videoVersionsMatch := videoVersionsRegex.FindStringSubmatch(jsonData)
	if len(videoVersionsMatch) < 2 {
		return nil
	}

	// Extract URLs from the video_versions section
	urlPattern := `"url":\s*"([^"]+)"`
	urlRegex := regexp.MustCompile(urlPattern)

	matches := urlRegex.FindAllStringSubmatch(videoVersionsMatch[1], -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}
