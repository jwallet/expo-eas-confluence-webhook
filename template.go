package main

import (
	"fmt"
	"strings"
)

func minify(template string) string {
	minifier := strings.NewReplacer("\n", "", "\t", "")
	return minifier.Replace(template)
}

func getPlatformText(platform string) string {
	if platform == "android" {
		return "Android"
	}
	return "iOS"
}

func getBuildTemplate(environment string, platform string, version string, sdk string, buildURL string, completedAt string, expiresAt string) string {
	platformText := getPlatformText(platform)
	template := fmt.Sprintf(`
	<table data-layout="default" ac:local-id="%v">
		<colgroup>
			<col style="width: 150.0px;" />
			<col style="width: 60.0px;" />
			<col style="width: 80.0px;" />
		</colgroup>
		<tbody>
			<tr>
				<td>
					<p><strong>%v</strong></p>
				</td>
				<td>
					<p style="text-align: center;"><strong>%v</strong></p>
				</td>
				<td>
					<p style="text-align: center;"><strong>SDK %v</strong></p>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<ac:structured-macro ac:name="iframe" ac:schema-version="1" data-layout="default">
						<ac:parameter ac:name="longdesc">Scan QR Code to install</ac:parameter>
						<ac:parameter ac:name="scrolling">no</ac:parameter>
						<ac:parameter ac:name="src"><ri:url ri:value="https://api.qrserver.com/v1/create-qr-code?size=200x200&amp;data=%v" /></ac:parameter>
						<ac:parameter ac:name="width">200</ac:parameter>
						<ac:parameter ac:name="frameborder">hide</ac:parameter>
						<ac:parameter ac:name="align">middle</ac:parameter>
						<ac:parameter ac:name="title">QR Code</ac:parameter>
						<ac:parameter ac:name="height">200</ac:parameter>
					</ac:structured-macro>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<p><a href="%v" data-card-appearance="inline">%v</a></p>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<p>Completed at: <strong>%v</strong></p>
				</td>
			</tr>
		</tbody>
	</table>`, environment+"-"+platform, platformText, version, sdk, buildURL, buildURL, buildURL, completedAt)
	return minify(template)
}

func getDefaultEnvironmentTemplate(environment string) string {
	android := getBuildTemplate(environment, "android", "1.0", "1.0", "http://httpstat.us/200", "2023-01-01T12:00:00.000Z", "2023-01-01T12:00:00.000Z")
	ios := getBuildTemplate(environment, "ios", "1.0", "1.0", "http://httpstat.us/200", "2023-01-01T12:00:00.000Z", "2023-01-01T12:00:00.000Z")
	template := fmt.Sprintf(`
	<ac:layout-section ac:type="two_equal" ac:breakout-mode="default">
		<ac:layout-cell>%v</ac:layout-cell>
		<ac:layout-cell>%v</ac:layout-cell>
	</ac:layout-section>`, android, ios)
	return minify(template)
}

func getDefaultTemplate() string {
	template := fmt.Sprintf(`
	<ac:layout>
		<ac:layout-section ac:type="fixed-width" ac:breakout-mode="default">
			<ac:layout-cell>
				<h2>Review App</h2>
			</ac:layout-cell>
		</ac:layout-section>
		%v
		<ac:layout-section ac:type="fixed-width" ac:breakout-mode="default">
			<ac:layout-cell>
				<h2>Continuous</h2>
			</ac:layout-cell>
		</ac:layout-section>
		%v
		<ac:layout-section ac:type="fixed-width" ac:breakout-mode="default">
			<ac:layout-cell>
				<h2>Integration</h2>
			</ac:layout-cell>
		</ac:layout-section>
		%v
		<ac:layout-section ac:type="fixed-width" ac:breakout-mode="default">
			<ac:layout-cell>
				<h2>Staging</h2>
			</ac:layout-cell>
		</ac:layout-section>
		%v
		<ac:layout-section ac:type="fixed-width" ac:breakout-mode="default">
			<ac:layout-cell>
				<h2>Production</h2>
			</ac:layout-cell>
		</ac:layout-section>
		%v
	</ac:layout>
`,
		getDefaultEnvironmentTemplate(string(review)),
		getDefaultEnvironmentTemplate(string(continuous)),
		getDefaultEnvironmentTemplate(string(staging)),
		getDefaultEnvironmentTemplate(string(integration)),
		getDefaultEnvironmentTemplate(string(production)))
	return minify(template)
}

// func getBuildTemplate(environment string, platform string, version string, sdk string, buildURL string, completedAt string, expiresAt string) string {
// 	platformText := getPlatformText(platform)
// 	template := fmt.Sprintf(`
// 	<table data-layout="default" ac:local-id="%v">
// 		<colgroup>
// 			<col style="width: 150.0px;" />
// 			<col style="width: 60.0px;" />
// 			<col style="width: 80.0px;" />
// 		</colgroup>
// 		<tbody>
// 			<tr>
// 				<td>
// 					<p><strong>%v</strong></p>
// 				</td>
// 				<td>
// 					<p style="text-align: center;"><strong>%v</strong></p>
// 				</td>
// 				<td>
// 					<p style="text-align: center;"><strong>SDK %v</strong></p>
// 				</td>
// 			</tr>
// 			<tr>
// 				<td colspan="3">
// 					<ac:structured-macro ac:name="iframe" ac:schema-version="1" data-layout="default">
// 						<ac:parameter ac:name="longdesc">Scan QR Code to install</ac:parameter>
// 						<ac:parameter ac:name="scrolling">no</ac:parameter>
// 						<ac:parameter ac:name="src"><ri:url ri:value="https://api.qrserver.com/v1/create-qr-code?size=200x200&amp;data=%v" /></ac:parameter>
// 						<ac:parameter ac:name="width">200</ac:parameter>
// 						<ac:parameter ac:name="frameborder">hide</ac:parameter>
// 						<ac:parameter ac:name="align">middle</ac:parameter>
// 						<ac:parameter ac:name="title">QR Code</ac:parameter>
// 						<ac:parameter ac:name="height">200</ac:parameter>
// 					</ac:structured-macro>
// 				</td>
// 			</tr>
// 			<tr>
// 				<td colspan="3">
// 					<p><a href="%v" data-card-appearance="inline">%v</a></p>
// 				</td>
// 			</tr>
// 			<tr>
// 				<td>
// 					<p>Completed at: <strong>%v</strong></p>
// 				</td>
// 				<td colspan="2">
// 					<p style="text-align: right;">Expires at: <strong>%v</strong></p>
// 				</td>
// 			</tr>
// 		</tbody>
// 	</table>`, environment+"-"+platform, platformText, version, sdk, buildURL, buildURL, buildURL, completedAt, expiresAt)
// 	return minify(template)
// }
