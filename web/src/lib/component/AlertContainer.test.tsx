import React, {ReactElement} from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import AlertContainer from './AlertContainer';
import {Icon, Snackbar} from '@material-ui/core';
import {runInAction} from 'mobx';

describe('AlertContainer', () => {
    const rootStore = new RootStore();
    const {alertStore} = rootStore;

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
        runInAction(() => alertStore.alerts.clear());
    });
    it('returns null if no alerts', () => {
        const component = shallow(<AlertContainer />);

        expect(component).toBeEmptyRender();
    });
    it('displays the first alert', () => {
        alertStore.addAlert('info', 'First alert');
        alertStore.addAlert('error', 'Second alert');

        const component = shallow(<AlertContainer />);

        const message = shallow(component.find(Snackbar).prop<ReactElement>('message'));
        expect(message).toHaveText('First alert');
    });
    it('displays the error icon for an error alert', () => {
        alertStore.addAlert('error', 'First alert');

        const component = shallow(<AlertContainer />);

        const message = shallow(component.find(Snackbar).prop<ReactElement>('message'));
        expect(message.find(Icon)).toHaveText('error_outline');
    });
    it('displays action button', () => {
        const onClick = jest.fn();
        alertStore.addAlert('info', 'First alert', {text: 'Do Something', onClick});
        const component = shallow(<AlertContainer />);

        const actionButton = component.find(Snackbar).prop<ReactElement<{onClick: () => void, children: unknown}>[]>('action')[0];

        expect(actionButton.props.children).toEqual('Do Something');
        actionButton.props.onClick();
        expect(onClick).toBeCalledTimes(1);
    });
    it('hides alert when close button is clicked', () => {
        alertStore.addAlert('error', 'First alert');
        const component = shallow(<AlertContainer />);
        const closeAction = component.find(Snackbar).prop<ReactElement<{onClick: () => void}>[]>('action')[1];

        closeAction.props.onClick();

        expect(alertStore.alerts[0].open).toBe(false);
    });
    it('removes alert after exit transition', () => {
        alertStore.addAlert('error', 'First alert');
        const component = shallow(<AlertContainer />);

        component.find(Snackbar).prop('onExited')?.({} as HTMLElement);

        expect(alertStore.alerts).toHaveLength(0);
    });
});
