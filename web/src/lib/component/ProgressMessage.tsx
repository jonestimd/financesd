import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';
import Typography from '@material-ui/core/Typography';

const ProgressMessage: React.FC = observer(() => {
    const rootStore = React.useContext(RootStoreContext);
    const message = rootStore.messageStore.progressMessage;
    return (
        <Typography className={'progress-overlay' + (message ? '' : ' hidden')}>
            <div className='progress-circular'>
                <div className='message'><p>{message}</p></div>
                <div className='progress-circular-wrapper'>
                    <div className='progress-circular-inner'>
                        <div className='progress-circular-left'>
                            <div className='progress-circular-spinner' />
                        </div>
                        <div className='progress-circular-right'>
                            <div className='progress-circular-spinner' />
                        </div>
                    </div>
                </div>
            </div>
        </Typography>
    );
});

export default ProgressMessage;
