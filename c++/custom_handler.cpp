#include <stdio.h>
#include <boost/lexical_cast.hpp>
#include "request.hpp"
#include "reply.hpp"
#include "mime_types.hpp"

int shell_cmd(const char *cmd, char *output, int maxlen)
{
    if (0 == cmd || 0 == output)
        return -1;

    memset(output, 0, maxlen);
    FILE *stream = popen(cmd, "r");
    return fread(output, sizeof(char), maxlen, stream);
}

const char* cmd_info = "/root/projects/ssctl/bin/ssctl info";
const char* cmd_reset = "/root/projects/ssctl/bin/ssctl reset";

bool handle_custom(const http::server2::request& req, http::server2::reply& rep)
{
    if (req.uri != "/~ss~/info" && req.uri != "/~ss~/reset") {
        return false;
    }
    char buf[512] = {0};
    if (req.uri[6] == 'i') {
        shell_cmd(cmd_info, buf, 512);
    } else {
        shell_cmd(cmd_reset, buf, 512);
    }

    rep.content.append(buf);
    rep.status = http::server2::reply::ok;
    rep.headers.resize(2);
    rep.headers[0].name = "Content-Length";
    rep.headers[0].value = boost::lexical_cast<std::string>(rep.content.size());
    rep.headers[1].name = "Content-Type";
    rep.headers[1].value = http::server2::mime_types::extension_to_type("html");

    return true;
}


