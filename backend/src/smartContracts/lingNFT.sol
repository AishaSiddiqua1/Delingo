// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract BadgeNFT is ERC721URIStorage, Ownable {
    uint256 public nextTokenId;

    constructor() ERC721("BadgeNFT", "BNFT") Ownable(msg.sender) {
        // Initialize token ID counter
        nextTokenId = 1;
    }

    // Function to mint a new NFT
    function mintBadge(
        address to,
        string memory name,
        string memory description,
        string memory criteria
    ) public onlyOwner {
        uint256 tokenId = nextTokenId;
        string memory tokenMetadata = createTokenURI(name, description, criteria);

        // Mint the NFT and set its metadata
        _mint(to, tokenId);
        _setTokenURI(tokenId, tokenMetadata);

        nextTokenId++;
    }

    // Create a token URI based on metadata
    function createTokenURI(
        string memory name,
        string memory description,
        string memory criteria
    ) public pure returns (string memory) {
        return
            string(
                abi.encodePacked(
                    "data:application/json;base64,",
                    encodeMetadata(name, description, criteria)
                )
            );
    }

    // Encode metadata into base64 JSON
    function encodeMetadata(
        string memory name,
        string memory description,
        string memory criteria
    ) public pure returns (string memory) {
        return
            string(
                abi.encodePacked(
                    '{"name":"',
                    name,
                    '", "description":"',
                    description,
                    '", "criteria":"',
                    criteria,
                    '"}'
                )
            );
    }

    // Set token URI (internal function)
    function _setTokenURI(uint256 tokenId, string memory newTokenURI) internal override {
        super._setTokenURI(tokenId, newTokenURI); // Call the parent function
    }
}
