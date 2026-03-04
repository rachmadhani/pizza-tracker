const statuses = ["Order Placed","Preparing","Baking","Ready"];
const currentStatus = "Order Placed";
const currentIndex = statuses.indexOf(currentStatus);
const progressPercentage = Math.max((currentIndex / (statuses.length - 1)) * 100, 0);
console.log(progressPercentage);
