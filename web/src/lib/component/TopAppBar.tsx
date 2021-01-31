import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import Icon from '@material-ui/core/Icon';
import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Menu from '@material-ui/core/Menu';

export interface ITopAppBarProps {
    title?: string;
    menuItems: React.ReactNode;
    children?: React.ReactNode;
}

const TopAppBar: React.FC<ITopAppBarProps> = ({title, menuItems, children}) => {
    const menuAnchor = React.useRef<HTMLDivElement>(null);
    const [showMenu, setShowMenu] = React.useState(false);
    const toggleMenu = React.useCallback(() => setShowMenu(!showMenu), [showMenu]);
    const hideMenu = React.useCallback(() => setShowMenu(false), []);
    return (
        <AppBar position='fixed'>
            <Toolbar variant='dense' disableGutters>
                <IconButton onClick={toggleMenu} color='inherit'><Icon>menu</Icon></IconButton>
                {title ? <Typography variant='h6'>{title}</Typography> : null}
                {children}
            </Toolbar>
            <div ref={menuAnchor} className='mdc-menu-surface--anchor' />
            <Menu open={showMenu} onClose={hideMenu} getContentAnchorEl={null} anchorEl={menuAnchor.current}>
                {menuItems}
            </Menu>
        </AppBar>
    );
};

export default TopAppBar;
