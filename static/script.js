// Game state
let gameState = null;

// DOM elements
const statusElement = document.getElementById('status');
const gameBoardElement = document.getElementById('game-board');
const cells = document.querySelectorAll('.cell');

// Initialize the game when page loads
document.addEventListener('DOMContentLoaded', function() {
    loadGameState();
});

// Load current game state from server
async function loadGameState() {
    try {
        const response = await fetch('/api/game');
        gameState = await response.json();
        updateUI();
    } catch (error) {
        console.error('Error loading game state:', error);
        statusElement.textContent = 'エラーが発生しました';
    }
}

// Make a move
async function makeMove(position) {
    if (gameState.gameOver || gameState.board[position] !== '') {
        return;
    }

    try {
        const response = await fetch('/api/move', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ position: position })
        });

        if (response.ok) {
            gameState = await response.json();
            updateUI();
        } else {
            console.error('Invalid move');
        }
    } catch (error) {
        console.error('Error making move:', error);
        statusElement.textContent = 'エラーが発生しました';
    }
}

// Start a new game
async function startNewGame() {
    try {
        const response = await fetch('/api/new-game', {
            method: 'POST'
        });

        if (response.ok) {
            gameState = await response.json();
            updateUI();
        }
    } catch (error) {
        console.error('Error starting new game:', error);
        statusElement.textContent = 'エラーが発生しました';
    }
}

// Update the UI based on current game state
function updateUI() {
    if (!gameState) return;

    // Update board
    cells.forEach((cell, index) => {
        cell.textContent = gameState.board[index];
        cell.setAttribute('data-symbol', gameState.board[index]);
        
        // Disable cells if game is over or cell is occupied
        if (gameState.gameOver || gameState.board[index] !== '') {
            cell.classList.add('disabled');
        } else {
            cell.classList.remove('disabled');
        }
    });

    // Update status
    if (gameState.gameOver) {
        if (gameState.winner === 'draw') {
            statusElement.textContent = '引き分けです！';
        } else {
            statusElement.textContent = `${gameState.winner}の勝利です！`;
        }
    } else {
        statusElement.textContent = `現在のターン: ${gameState.currentTurn}`;
    }

    // Add winner animation if there's a winner
    if (gameState.gameOver && gameState.winner !== 'draw') {
        highlightWinningCells();
    }
}

// Highlight winning cells with animation
function highlightWinningCells() {
    const winPatterns = [
        [0, 1, 2], [3, 4, 5], [6, 7, 8], // rows
        [0, 3, 6], [1, 4, 7], [2, 5, 8], // columns
        [0, 4, 8], [2, 4, 6] // diagonals
    ];

    for (const pattern of winPatterns) {
        if (gameState.board[pattern[0]] !== '' &&
            gameState.board[pattern[0]] === gameState.board[pattern[1]] &&
            gameState.board[pattern[1]] === gameState.board[pattern[2]]) {
            
            pattern.forEach(index => {
                cells[index].classList.add('winner-cell');
            });
            break;
        }
    }
}

// Add click sound effect (optional enhancement)
function playClickSound() {
    // This could be enhanced with actual sound files
    // For now, we'll just use a simple visual feedback
}

// Keyboard support
document.addEventListener('keydown', function(event) {
    if (gameState && gameState.gameOver) return;

    const keyMap = {
        '1': 6, '2': 7, '3': 8,
        '4': 3, '5': 4, '6': 5,
        '7': 0, '8': 1, '9': 2
    };

    if (keyMap.hasOwnProperty(event.key)) {
        event.preventDefault();
        makeMove(keyMap[event.key]);
    }

    // New game with 'n' key
    if (event.key.toLowerCase() === 'n') {
        event.preventDefault();
        startNewGame();
    }
});