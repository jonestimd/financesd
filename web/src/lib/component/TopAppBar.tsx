import React from 'react';
import IconButton from '@material/react-icon-button';
import MaterialIcon from '@material/react-material-icon';
import MDCAppBar, {TopAppBarRow, TopAppBarSection, TopAppBarTitle} from '@material/react-top-app-bar';
import Menu, {MenuList} from '@material/react-menu';

export interface ITopAppBarProps {
    title: string;
    menuItems: React.ReactNode;
}

const TopAppBar: React.FC<ITopAppBarProps> = ({title, menuItems}) => {
    const menuAnchor = React.useRef(null);
    const [showMenu, setShowMenu] = React.useState(false);
    const toggleMenu = React.useCallback(() => setShowMenu(!showMenu), [showMenu]);
    const hideMenu = React.useCallback(() => setShowMenu(false), []);
    return (
        <MDCAppBar fixed dense>
            <TopAppBarRow>
                <TopAppBarSection>
                    <IconButton onClick={toggleMenu}><MaterialIcon icon='menu' /></IconButton>
                    <TopAppBarTitle>{title}</TopAppBarTitle>
                </TopAppBarSection>
            </TopAppBarRow>
            <div ref={menuAnchor} className='mdc-menu-surface--anchor'>
                <Menu open={showMenu} onClose={hideMenu} onSelected={hideMenu} anchorElement={menuAnchor.current}>
                    <MenuList>{menuItems}</MenuList>
                </Menu>
            </div>
        </MDCAppBar>
    );
};

export default TopAppBar;