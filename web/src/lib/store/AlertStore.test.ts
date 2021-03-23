import {runInAction} from 'mobx';
import AlertStore from './AlertStore';

describe('AlertStore', () => {
    const alertStore = new AlertStore();

    beforeEach(() => {
        runInAction(() => alertStore.alerts.clear());
    });
    describe('addAlert', () => {
        it('appends an info alert', () => {
            const onClick = () => console.info('here');

            alertStore.addAlert('info', 'First Alert', {text: 'Do Something', onClick});

            expect(alertStore.alerts).toHaveLength(1);
            expect(alertStore.alerts[0]).toEqual(expect.objectContaining({type: 'info', message: 'First Alert', open: true}));
            expect(alertStore.alerts[0].action).toEqual({text: 'Do Something', onClick});
            expect(alertStore.alerts[0].timerId).toBeDefined();
        });
        it('appends an error alert', () => {
            alertStore.addAlert('error', 'Next Alert');

            expect(alertStore.alerts).toHaveLength(1);
            expect(alertStore.alerts[0]).toEqual(expect.objectContaining({type: 'error', message: 'Next Alert', open: true}));
            expect(alertStore.alerts[0].timerId).not.toBeDefined();
        });
    });
    describe('hideAlert', () => {
        it('sets open to false', () => {
            alertStore.addAlert('error', 'First Alert');

            alertStore.hideAlert(alertStore.alerts[0].id);

            expect(alertStore.alerts[0].open).toBe(false);
        });
        it('ignores unknown id', () => {
            alertStore.hideAlert(-1);
        });
    });
    describe('removeAlert', () => {
        it('removes an alert', () => {
            alertStore.addAlert('error', 'First Alert');

            alertStore.removeAlert(alertStore.alerts[0].id);

            expect(alertStore.alerts).toHaveLength(0);
        });
        it('clears timer', () => {
            jest.spyOn(window, 'setTimeout');
            jest.spyOn(window, 'clearTimeout');
            alertStore.addAlert('info', 'First Alert');

            alertStore.removeAlert(alertStore.alerts[0].id);

            expect(window.clearTimeout).toBeCalledTimes(1);
        });
        it('ignores unknown id', () => {
            alertStore.removeAlert(-1);
        });
    });
});
