import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { baseHref } from '@/portainer/helpers/pathHelper';

export default class KubectlShellController {
  /* @ngInject */
  constructor(TerminalWindow, $window, $async, EndpointProvider, LocalStorage, Notifications) {
    this.$async = $async;
    this.$window = $window;
    this.TerminalWindow = TerminalWindow;
    this.EndpointProvider = EndpointProvider;
    this.LocalStorage = LocalStorage;
    this.Notifications = Notifications;

    $window.onbeforeunload = () => {
      if (this.state.shell.connected) {
        return '';
      }
    };
  }

  disconnect() {
    if (this.state.shell.connected) {
      this.state.shell.connected = false;
      this.state.icon = 'fas fa-window-minimize';
      this.state.shell.socket.close();
      this.state.shell.term.dispose();
      this.TerminalWindow.terminalclose();
      this.$window.onresize = null;
    }
  }

  screenClear() {
    this.state.shell.term.clear();
  }

  miniRestore() {
    if (this.state.css === 'mini') {
      this.state.css = 'normal';
      this.state.icon = 'fas fa-window-minimize';
      this.TerminalWindow.terminalopen();
    } else {
      this.state.css = 'mini';
      this.state.icon = 'fas fa-window-restore';
      this.TerminalWindow.terminalclose();
    }
  }

  configureSocketAndTerminal(socket, term, fitAddon) {
    socket.onopen = () => {
      const terminal_container = document.getElementById('terminal-container');
      term.loadAddon(fitAddon);
      term.open(terminal_container);
      term.setOption('cursorBlink', true);
      term.focus();
      fitAddon.fit();
      term.writeln('#Run kubectl commands inside here');
      term.writeln('#e.g. kubectl get all');
      term.writeln('');
    };

    term.onData(function (data) {
      socket.send(data);
    });

    socket.onmessage = (msg) => {
      term.write(msg.data);
    };

    socket.onerror = (err) => {
      this.disconnect();
      if (err.target.readyState !== WebSocket.CLOSED) {
        this.Notifications.error('Failure', err, 'Websocket connection error');
      }
    };

    this.$window.onresize = () => {
      this.TerminalWindow.terminalresize();
    };

    socket.onclose = this.disconnect.bind(this);

    this.state.shell.connected = true;
  }

  connectConsole() {
    this.TerminalWindow.terminalopen();
    this.state.css = 'normal';

    const params = {
      token: this.LocalStorage.getJWT(),
      endpointId: this.EndpointProvider.endpointID(),
    };

    const wsProtocol = this.$window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    const path = baseHref() + 'api/websocket/kubernetes-shell';
    const base = path.startsWith('http') ? path.replace(/^https?:\/\//i, '') : window.location.host + path;

    const queryParams = Object.entries(params)
      .map(([k, v]) => `${k}=${v}`)
      .join('&');

    const url = `${wsProtocol}${base}?${queryParams}`;
    this.state.shell.socket = new WebSocket(url);
    this.state.shell.term = new Terminal({ cursorBlink: true });
    this.state.shell.fitAddon = new FitAddon();

    this.configureSocketAndTerminal(this.state.shell.socket, this.state.shell.term, this.state.shell.fitAddon);
  }

  $onInit() {
    return this.$async(async () => {
      this.state = {
        css: 'normal',
        icon: 'fa-window-minimize',
        shell: {
          connected: false,
          socket: null,
          term: null,
        },
      };
    });
  }

  $onDestroy() {
    if (this.state.shell.connected) {
      this.disconnect();
      this.$window.onresize = null;
    }
  }
}
