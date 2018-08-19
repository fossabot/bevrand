from flask import Blueprint, jsonify, request
from flasgger import swag_from
from api.services import data_validator
from api.setup import  FLASK_TRACER
from api.error_handler.error_model import InvalidUsage
import opentracing
from api.services.frontpage_service import FrontPageService


front_page_blueprint = Blueprint('front_page', __name__,)


@front_page_blueprint.route('/ping', methods=['GET'])
def ping_pong():
    return jsonify({
        'status': 'success',
        'message': 'pong!'
    })


@front_page_blueprint.route('', methods=['GET'])
@swag_from('../swagger/public_get.yml')
def front_page_all_lists():
    parent_span = create_parent_trace()
    with opentracing.tracer.start_span('playlist_frontpage', child_of=parent_span) as span:
        service = FrontPageService()
        result = service.retrieve_all_front_page_lists()
        return jsonify({"result": result}), 200


@front_page_blueprint.route('/<play_list_name>', methods=['GET'])
@swag_from('../swagger/public_get_playlist.yml')
def front_page_list(play_list_name):
    parent_span = create_parent_trace()
    with opentracing.tracer.start_span('playlist_frontpage', child_of=parent_span) as span:
        data_validator.validate_play_list(play_list_name)
        service = FrontPageService()
        result = service.retrieve_front_page_list(play_list_name)
        return jsonify({"result": result.__dict__}), 200


@front_page_blueprint.errorhandler(InvalidUsage)
def handle_invalid_usage(error):
    response = jsonify(error.to_dict())
    response.status_code = error.status_code
    parent_span = create_parent_trace()
    with opentracing.tracer.start_span('error', child_of=parent_span) as span:
        span.log_kv({"status_code": error.status_code, "error": error.message})
        if error.meta is not None:
            span.log_kv({"meta": error.meta})
    return response


@FLASK_TRACER.trace()
def create_parent_trace():
    parent_span = FLASK_TRACER.get_span(request)
    parent_span.set_tag('http.url', request.base_url)
    parent_span.set_tag('http.method', request.method)
    parent_span.set_tag('body', request.json)
    return parent_span

