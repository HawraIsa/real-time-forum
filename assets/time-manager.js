// class TimeManager {
//     constructor() {
//         this.interval = null;
//         this.UPDATE_INTERVAL = 60000; // 1 minute
//         this.boundUpdate = this.updateAllTimes.bind(this);
//     }

//     start() {
//         if (this.interval) return;
        
//         // Initial update
//         this.updateAllTimes();
        
//         // Set up independent interval
//         this.interval = setInterval(this.boundUpdate, this.UPDATE_INTERVAL);
        
//         // Watch for DOM changes
//         this.setupObserver();
//     }

//     stop() {
//         if (this.interval) {
//             clearInterval(this.interval);
//             this.interval = null;
//         }
//         if (this.observer) {
//             this.observer.disconnect();
//         }
//     }

//     setupObserver() {
//         this.observer = new MutationObserver((mutations) => {
//             let shouldUpdate = false;
//             mutations.forEach((mutation) => {
//                 if (mutation.addedNodes.length) {
//                     const hasTimeElements = Array.from(mutation.addedNodes).some(node => 
//                         node.nodeType === 1 && (
//                             node.querySelector('[data-timestamp]') || 
//                             node.hasAttribute('data-timestamp')
//                         )
//                     );
//                     if (hasTimeElements) {
//                         shouldUpdate = true;
//                     }
//                 }
//             });
//             if (shouldUpdate) {
//                 this.updateAllTimes();
//             }
//         });

//         this.observer.observe(document.body, {
//             childList: true,
//             subtree: true
//         });
//     }

//     updateAllTimes() {
//         document.querySelectorAll('[data-timestamp]').forEach(element => {
//             const timestamp = parseInt(element.dataset.timestamp);
//             if (timestamp) {
//                 element.textContent = this.formatTime(timestamp);
//             }
//         });
//     }

//     formatTime(timestamp) {
//         const now = new Date();
//         const messageDate = new Date(timestamp);
//         const diff = now - messageDate;
//         const isToday = now.toDateString() === messageDate.toDateString();
//         const isThisYear = now.getFullYear() === messageDate.getFullYear();

//         if (diff < 60000) return 'now';
//         if (diff < 3600000) return `${Math.floor(diff / 60000)}m`;
//         if (isToday) return messageDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
//         if (diff < 86400000 * 2) return 'Yesterday';
//         if (diff < 86400000 * 7) return messageDate.toLocaleDateString([], { weekday: 'short' });
//         if (isThisYear) return messageDate.toLocaleDateString([], { month: 'short', day: 'numeric' });
//         return messageDate.toLocaleDateString([], { year: 'numeric', month: 'short', day: 'numeric' });
//     }
// }

// // Single instance
// const timeManager = new TimeManager();

// // Independent lifecycle
// document.addEventListener('DOMContentLoaded', () => timeManager.start());
// window.addEventListener('unload', () => timeManager.stop());