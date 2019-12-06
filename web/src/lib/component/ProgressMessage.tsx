import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';

const ProgressMessage: React.FC<{}> = observer(() => {
    const rootStore = React.useContext(RootStoreContext);
    const message = rootStore.progressMessage;
    return message ? <div className='progress-indicator'>{message}</div> : null;
});

export default ProgressMessage;