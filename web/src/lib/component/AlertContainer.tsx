import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../store/RootStore';
import {Button, Icon, IconButton, Slide, Snackbar, SnackbarOrigin} from '@material-ui/core';

const anchor: SnackbarOrigin = {vertical: 'bottom', horizontal: 'center'};

const AlertContainer: React.FC = observer(() => {
    const {alertStore} = React.useContext(RootStoreContext);
    const [alert] = alertStore.alerts;
    if (alert) {
        const onClose = () => alertStore.hideAlert(alert.id);
        const message = alert.type === 'error'
            ? <span><Icon fontSize='small'>error_outline</Icon>{alert.message}</span>
            : <span>{alert.message}</span>;
        const action = alert.action
            ? <Button key='action' color='primary' onClick={alert.action.onClick}>{alert.action.text}</Button>
            : null;
        const closeAction = <IconButton key='close' color='inherit' size='small' onClick={onClose}><Icon fontSize='small'>close</Icon></IconButton>;
        return (
            <Snackbar key={alert.id} open={alert.open} anchorOrigin={anchor} message={message} action={[action, closeAction]}
                onExited={() => alertStore.removeAlert(alert.id)} TransitionComponent={Slide} />
        );
    }
    return null;
});

export default AlertContainer;
