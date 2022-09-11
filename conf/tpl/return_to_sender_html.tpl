<html>
<head>
<title>{TITLE}</title>
</head>
<body>
<p style="font-size:14px; font-family:verdana; margin:0 0 10px;
line-height:22px;
text-align:left">很抱歉您发送的邮件被退回，以下是该邮件的相关信息：</p>

<table style="width:90%; border-collapse:collapse; margin:0 0 15px" cellpadding="0" cellspacing="0" border="0">
	<tr>
		<td style="width:120px; font-family:verdana; font-size:14px; font-weight:bold; text-align:center; background:#E0ECF9; border:1px solid #C1D9F3">被退回邮件</td>
		<td style="border:1px solid #C1D9F3; font-family:verdana; padding:5px 15px; line-height:20px; font-size:14px;">主 题：{SEND_SUBJECT} <br/>时 间：{TIME} </td>
	</tr>
</table>

<table style="width:90%; border-collapse:collapse; margin:0 0 10px" cellpadding="0" cellspacing="0" border="0">
	<tr>
		<td colspan="2" style="font-size:14px;font-weight:bold;padding:5px 10px;border:1px solid #C1D9F3; background:#C1D9F3">无法发送到 {ERR_TO_MAIL}</td>
	</tr>
	<tr>
		<td style="width:120px; font-size:14px;text-align:center; background:#EFF5FB; border:1px solid #C1D9F3">退信原因</td>
		<td style="border:1px solid #C1D9F3; font-family:verdana; padding:5px 15px; line-height:20px; font-size:14px;"> 未知原因 <br /> {ERR_MSG} </td>
	</tr>
	<tr>
		<td style="width:120px;font-size:14px;text-align:center; background:#EFF5FB; border:1px solid #C1D9F3">解决方案</td>
		<td style="border:1px solid #C1D9F3; font-family:verdana; padding:5px 15px; line-height:20px; font-size:14px;"> 请联系您的收件人（{ERR_TO_MAIL}）或其所在服务商了解退信原因。</td></tr>

</table>

<p style="font-size:14px; font-family:verdana; margin:0 0 10px; line-height:22px; text-align:left">此外，您还可以 <a href="https://github.com/phper95/mail-server/wiki/%E9%80%80%E4%BF%A1" target="_blank">点击这里</a> 获取更多关于退信的帮助信息。</p>

<hr>
<p style="font-size:14px; font-family:verdana; margin:0 0 10px; line-height:22px; text-align:left">某些云服务器的25外接端口，有可能被封，需要和云服务官方联系开通！</p>
</body>
</html>
