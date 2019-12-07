import React from 'react';
import {RootStoreContext} from '../store/RootStore';
import {observer} from 'mobx-react-lite';

const ProgressMessage: React.FC<{}> = observer(() => {
    const rootStore = React.useContext(RootStoreContext);
    const message = rootStore.progressMessage;
    return (
        <div className={'progress-overlay' + (message ? '' : ' hidden')}>
            <div className="progress-circular">
                <div className='message'>{message}</div>
                <div className="progress-circular-wrapper">
                    <div className="progress-circular-inner">
                        <div className="progress-circular-left">
                            <div className="progress-circular-spinner"></div>
                        </div>
                        <div className="progress-circular-right">
                            <div className="progress-circular-spinner"></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
});

export default ProgressMessage;