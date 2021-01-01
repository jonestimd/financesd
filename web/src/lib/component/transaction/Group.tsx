import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';

const Group: React.FC<{id?: number}> = observer(({id}) => {
    const {groupStore} = React.useContext(RootStoreContext);
    return typeof id === 'number' ? <span className='group'>{groupStore.getGroup(id).name}</span> : null;
});

export default Group;
