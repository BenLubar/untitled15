var command = function() {
	var audio = {};

	return {
		'': function(u, n) {
			audio[u].next.push(n);
		},
		'connected': function(u) {
			var context = new (window.AudioContext || window.webkitAudioContext)();
			var input = context.createOscillator();
			var node = (context.createScriptProcessor ||
			            context.createJavaScriptNode).call(context, 4096, 1, 1);
			node.next = [];
			node.onaudioprocess = function(e) {
				var buf = this.next.shift();
				if (buf) {
					var data = e.outputBuffer.getChannelData(0);
					for (var i = 0; i < data.length; i++) {
						data[i] = buf[i];
					}
				}
			};
			input.connect(node);
			node.input = input;
			node.connect(context.destination);
			audio[u] = node;
		},
		'disconnected': function(u) {
			audio[u].input.disconnect();
			audio[u].disconnect();
			delete audio[u];
		}
	};
}();

var net = function() {
	var queue = [];
	var backoff = 0;
	var ws;
	function reset() {
		ws = new WebSocket('ws://' + location.host + '/sock');
		ws.onmessage = function(e) {
			var data = JSON.parse(e.data);
			if (data.Special) {
				command[data.Special](data.User);
			}
			if (data.Audio) {
				var audio = new Uint8Array(data.Audio.length);
				for (var i = 0; i < audio.length; i++) {
					audio[i] = data.Audio.charCodeAt(i);
				}
				command[''](data.User, new Float32Array(audio.buffer));
				data.Audio = true;
			}
			console.log('net: message:', data);
		};
		ws.onopen = function(e) {
			console.log('net: connection opened');
			queue.forEach(function(p) {
				ws.send(p);
			});
			queue = null;
		};
		ws.onclose = function(e) {
			console.log('net: connection closed:', e.code, e.reason, 'clean: ' + e.wasClean);
			if (!queue) {
				queue = [];
			}
			setTimeout(reset, backoff += 1000);
		};
		ws.onerror = function(e) {
			console.log('net: unknown error');
		};
	}

	reset();

	return {
		send: function(data) {
			data = JSON.stringify(data);
			if (queue) {
				queue.push(data);
			} else {
				ws.send(data);
			}
		}
	};
}();

(navigator.getUserMedia ||
 navigator.webkitGetUserMedia ||
 navigator.mozGetUserMedia ||
 navigator.msGetUserMedia ||
 function(opt, success, failure) {
	 failure('getUserMedia is not supported by this browser');
 }).call(navigator,
	{video: false, audio: true},
	function(stream) {
		console.log('stream', stream);
		var context = new (window.AudioContext || window.webkitAudioContext)();
		var input = context.createMediaStreamSource(stream);
		var node = (context.createScriptProcessor ||
		            context.createJavaScriptNode).call(context, 4096, 1, 1);
		node.onaudioprocess = function(e) {
			net.send({'Audio': String.fromCharCode.apply(String, new Uint8Array(e.inputBuffer.getChannelData(0).buffer))});
		};
		input.connect(node);
		node.connect(context.destination);
	}, function(err) {
		console.log('error', err);
	});
