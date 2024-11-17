// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract LingToken is ERC20, Ownable {
    uint256 public constant MAX_SUPPLY = 100000 * 10**18; // Fixed total supply
    uint256 public lastBurnTimestamp; // Tracks the last time tokens were burned

    constructor() ERC20("Ling", "LING") {
        _mint(msg.sender, MAX_SUPPLY); // Mint the full supply to the owner
    }

    // Burn function for users to burn their own tokens
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }

    // Automatic 20% burn function
    function burnCirculatingSupply() external onlyOwner {
        require(
            block.timestamp >= lastBurnTimestamp + 365 days,
            "Burn can only occur once per year"
        );

        uint256 circulatingSupply = totalSupply();
        uint256 burnAmount = (circulatingSupply * 20) / 100; // 20% of the circulating supply

        _burn(owner(), burnAmount); // Burn from owner's balance
        lastBurnTimestamp = block.timestamp; // Update the last burn timestamp
    }
}
