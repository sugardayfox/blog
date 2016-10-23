Id: 1189
Title: Results of tweaking compiler flags before 0.9 release:
Tags: sumatra,optimization,programming
Date: 2008-08-11T04:55:37-07:00
Format: Markdown
--------------
Results of tweaking compiler flags before 0.9 release:


** Original build

08/10/2008  11:33 AM         2,128,896 pdftool.exe

08/10/2008  11:33 AM         2,784,768 SumatraPDF.exe


** With /Os:

08/10/2008  11:38 AM         2,062,848 pdftool.exe

08/10/2008  11:38 AM         2,696,192 SumatraPDF.exe


** With /Og and without /RTCs and /RTCu

08/10/2008  11:41 AM         1,796,608 pdftool.exe

08/10/2008  11:41 AM         2,349,568 SumatraPDF.exe


** Withoug /Og

08/10/2008  11:43 AM         1,994,240 pdftool.exe

08/10/2008  11:43 AM         2,588,160 SumatraPDF.exe


** Using /Ox instead of /Og

08/10/2008  11:44 AM         1,789,952 pdftool.exe

08/10/2008  11:45 AM         2,345,472 SumatraPDF.exe


** Adding /Gy

08/10/2008  11:47 AM         1,725,952 pdftool.exe

08/10/2008  11:47 AM         2,222,080 SumatraPDF.exe


** Adding /GL and /LTCG

08/10/2008  11:50 AM         1,718,784 pdftool.exe

08/10/2008  11:50 AM         2,212,864 SumatraPDF.exe


