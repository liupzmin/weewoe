package tmpl

var tm = `
<html>
<head>
    <title>{{- .Title }}</title>
</head>
<body>
    <style type="text/css">
        body {
            font: bold 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: White;
        }

        table.tdiff {
            border-collapse: collapse;
        }

        th {
            font: bold 14pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: white;
            padding: 8px;
        }

        th.title {
            font: bold 34pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: white;
            border: 1px solid black;
        }

        td.summary {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: white;
            vertical-align: top;
            border-left: 1px solid black;
            border-bottom: 1px solid black;
        }

        #date,
        #checker,
        #attention {
            border-right: 1px solid black
        }

        tr.awrc {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: #F0F8FF;
            vertical-align: top;
        }

        tr.orangered {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: red;
            vertical-align: top;
        }

        tr.yellow {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: yellow;
            vertical-align: top;
        }

        tr.awrnc {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: black;
            background: White;
            vertical-align: top;
        }

        td.red {
            font: 10pt Arial, Helvetica, Geneva, sans-serif;
            color: rgb(250, 250, 250);
            background: red;
            vertical-align: top;
        }

        td.littletitle {
            text-align: center;
    		background: white;
    		font-family: emoji;
    		font-size: 12pt;
    		font-weight: bold;
    		display: table-cell;
    		vertical-align: inherit;
    		border-left: 2px solid skyblue;
    		color: #909399;
        }

        span.green {
            color: #1ad01a;
			font-weight: 600;
        }
		
		span.orange {
            color: orange;
			font-weight: 600;
        }

        .gblock {
            border: 1px solid #cedced;
            width: 18px;
            height: 18px;
            background: #1ad01a;
            float: left;
            margin-right: 4px;
            margin-left: 4px;
            font-size: 12px;
            text-align: center;
            line-height: 18px;
        }

        .yblock {
            border: 1px solid #cedced;
            width: 18px;
            height: 18px;
            background: yellow;
            float: left;
            margin-right: 4px;
            margin-left: 4px;
            font-size: 12px;
            text-align: center;
            line-height: 18px;
        }

        .rblock {
            border: 1px solid #cedced;
            width: 18px;
            height: 18px;
            background: red;
            float: left;
            margin-right: 4px;
            margin-left: 4px;
            font-size: 12px;
            text-align: center;
            line-height: 18px;
        }

        .oblock {
            border: 1px solid #cedced;
            width: 18px;
            height: 18px;
            background: orangered;
            float: left;
            margin-right: 4px;
            margin-left: 4px;
            font-size: 12px;
            text-align: center;
            line-height: 18px;
        }

        .tfloat {
            height: 20px;
            float: left;
            font-size: 12px;
            text-align: center;
            line-height: 20px;
        }
    </style>
    </head>

    <body bgcolor="#ffffff" style="font-size:15px" text="#000000">
        <table cellpadding="2" cellspacing="0" class="tdiff" width="850">
            <tbody>
                <tr>
                    <th class="title" colspan="6">{{ .ReportName }}</th>
                </tr>
                <tr>
                    <td class="summary" colspan="2">报告日期</td>
                    <td class="summary" colspan="4" id="date">{{ .ReportDate }}</td>
                </tr>
                <tr>
                    <td class="summary" colspan="2">巡检单位</td>
                    <td class="summary" colspan="4" id="checker">{{ .Reporter }}</td>
                </tr>
                <tr>
                    <td class="summary" colspan="2">注意事项</td>
                    <td class="summary" colspan="4" id="attention">
                        <div class="gblock">{{ .NormalCount }}</div>
                        <div class="tfloat">正常</div>
                        <div class="yblock">{{ .PlanStopCount }}</div>
                        <div class="tfloat">计划中停止</div>
                        <div class="oblock">{{ .PortExceptionCount }}</div>
                        <div class="tfloat">端口异常</div>
                        <div class="rblock">{{ .PIDExceptionCount }}</div>
                        <div class="tfloat">进程异常</div>
                    </td>
                </tr>
            </tbody>
        </table>
        <table class="tdiff" id="program" width="850">
            <tbody>
				{{- with .Groups }}
				{{- range . }}
					<tr>
                    	<th colspan="6">{{ .Name }}</th>
                	</tr>
					<tr>
                        <td class="littletitle">序号</td>
                        <td class="littletitle">程序名称</td>
                        <td class="littletitle">主机</td>
                        <td class="littletitle">端口</td>
                        <td class="littletitle">启动时间</td>
                        <td class="littletitle">进程状态</td>
                	</tr>
					{{- range $i, $p := .Processes }}
						{{- $class := "" }}
						{{- if isEven $i}}
							{{- $class = "awrc" }}
						{{- else }}
							{{- $class = "awrnc" }}
						{{- end }}
						{{- if eq $p.State 0 }}
							{{- $class = "orangered" }}
						{{- end }}
						{{- if $p.Suspend }}
							{{- $class = "yellow" }}
						{{- end }}
						<tr class="{{ $class }}">
                    		<td>{{ inc $i }}</td>
                    		<td>{{ $p.Name }}</td>
                    		<td>{{ $p.Host }}</td>
                    		<td>
								{{- with $p.Ports }}
								{{- range . }}
									{{- if eq .State 0 }}
										<span class="orange" >{{ .Number }}</span>
									{{- else }}
										<span class="green">{{ .Number }}</span>
									{{- end }}
								{{- end }}
								{{- end }}
							</td>
                    		<td>{{ $p.StartTime }}</td>
                    		<td>
								{{- if eq $p.State 0 }}
									<span class="orange" >Down</span>
								{{- else }}
									<span class="green">Up</span>
								{{- end }}
							</td>
                		</tr>
					{{- end }}
				{{- end }}
				{{- end }}
            </tbody>
        </table>
    </body>
</html>
`
