#ifndef _CUSTOM_HANDLER_H_
#define _CUSTOM_HANDLER_H_

#include "request.hpp"
#include "reply.hpp"

bool handle_custom(const http::server2::request& req, http::server2::reply& rep);

#endif
