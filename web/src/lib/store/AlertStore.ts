import {action, makeObservable, observable} from 'mobx';

export type AlertType = 'error' | 'info';

export interface IAlertAction {
    text: string;
    onClick: () => Promise<void> | void;
}

export class IAlert {
    readonly id = nextId++;
    @observable open = true;
    timerId?: number;

    constructor(
        public readonly type: AlertType,
        public readonly message: string,
        public readonly action?: IAlertAction) {
        makeObservable(this);
    }
}

let nextId = 1;

const infoTimeout = 3000;

export default class AlertStore {
    readonly alerts = observable.array<IAlert>([], {deep: false});

    constructor() {
        makeObservable(this);
    }

    @action
    addAlert(type: AlertType, message: string, action?: IAlertAction) {
        const alert = new IAlert(type, message, action);
        if (type === 'info') {
            alert.timerId = window.setTimeout(() => this.hideAlert(alert.id), infoTimeout);
        }
        this.alerts.push(alert);
    }

    @action
    hideAlert(id: number) {
        const alert = this.alerts.find((a) => a.id === id);
        if (alert) alert.open = false;
    }

    @action
    removeAlert(id: number) {
        const alert = this.alerts.find((a) => a.id === id);
        if (alert) {
            this.alerts.remove(alert);
            if (alert.timerId) window.clearTimeout(alert.timerId);
        }
    }
}
