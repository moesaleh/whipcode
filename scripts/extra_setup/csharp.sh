wget https://www.nuget.org/api/v2/package/Microsoft.Net.Compilers.Toolset/4.11.0 -O /tmp/Microsoft.Net.Compilers.Toolset.nupkg
unzip /tmp/Microsoft.Net.Compilers.Toolset.nupkg -d /tmp/Microsoft.Net.Compilers.Toolset
mv /tmp/Microsoft.Net.Compilers.Toolset/tasks/net472/* /usr/lib/mono/4.5/
rm -rf /tmp/Microsoft.Net.Compilers.Toolset

echo -e "
global using global::System;
global using global::System.Collections.Generic;
global using global::System.IO;
global using global::System.Linq;
global using global::System.Net.Http;
global using global::System.Threading;
global using global::System.Threading.Tasks;
" > /tmp/GlobalUsings.cs
