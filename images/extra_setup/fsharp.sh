wget https://www.nuget.org/api/v2/package/FSharp.Compiler.Tools/10.2.3 -O /tmp/FSharp.Compiler.Tools.nupkg
unzip /tmp/FSharp.Compiler.Tools.nupkg -d /tmp/FSharp.Compiler.Tools
mkdir -p /usr/local/fsharp/ && mv /tmp/FSharp.Compiler.Tools/tools/* /usr/local/fsharp/
rm -rf /tmp/FSharp.Compiler.Tools
rm -f /tmp/FSharp.Compiler.Tools.nupkg
